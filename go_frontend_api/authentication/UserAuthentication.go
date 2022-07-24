package authentication

import (
	"database/sql"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"bl.com/api/encryption"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gomodule/redigo/redis"

	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	_ "golang.org/x/oauth2"
	_ "golang.org/x/oauth2/google"
	_ "google.golang.org/api/urlshortener/v1"
)

var Db *sql.DB

var jwtKey = []byte("banana")

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

//Takes in one parameter called idToken
func Googlesignin(w http.ResponseWriter, r *http.Request) {
	//Validate the authenticity of idToken

	idToken := r.Header.Get("googleToken")
	//idToken := r.URL.Query().Get("idToken")

	googleClaimsObj, err := ValidateGoogleJWT(idToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Error message: " + err.Error())
		return
	}
	email := googleClaimsObj.Email
	username := googleClaimsObj.FirstName

	//Check if user email in database
	result := Db.QueryRow("select username from userinfo where email=?", email)
	if result == nil {
		//User has not been found, add them to database
		if _, err = Db.Query("insert into userinfo (username, email) values(?, ?)", username, email); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error message: " + err.Error())
			return
		}
		//Persist a JWT to the client
		sendJWT(w, email)
	}
	//User has been found, assign them a JWT
	sendJWT(w, email)
}

func sendJWT(w http.ResponseWriter, email string) {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Send response back with the private JWT Token String

	sendResponse(w, tokenString)

}

//REMEMBER, WITH VOLLEY, SERVER CAN ACCEPT HEADERS IN REQUESTS, BUT HAVE TO RESPOND WITH A RESPONSE, and NOT A HEADER,
// BECAUSE VOLLEY CANT ACCESS HEADERS FROM RECEIVING RESPONSES, IT CAN ONLY MODIFY HEADERS BEFORE SENDING THEM
func JWTAuthTester(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the request's header, which come with every request

	//Get the token header that is sent
	header := r.Header.Get("privateToken")

	if header == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Get the JWT string from the header
	tknStr := header

	fmt.Println("Private JWT Token String: " + tknStr)

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Invalid signature")
			return
		}
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		fmt.Println("Token not valid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user, along with their
	// email given in the token
	w.Header().Set("Content-Type", "application/json")

	sendResponse(w, "Welcome "+claims.Email+"!")

}

//Method to encrpyt and send responses
func sendResponse(w http.ResponseWriter, response string) {

	// To encrypt the StringToEncrypt
	encText, err := encryption.Encrypt(response, encryption.MySecret)
	if err != nil {
		fmt.Println("error encrypting your classified text: ", err)
		fmt.Fprint(w, "Error encrypting the request, please try again in a couple seconds")
	}
	fmt.Fprint(w, encText)

}

//To log all requests to the API
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		// Serve the request
		next.ServeHTTP(w, r)
	})
}

//Refresh token endpoint - to be called periodically from the frontend client
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	//Get the token header that is sent
	header := r.Header.Get("privateToken")

	if header == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	// Get the JWT string from the header
	tknStr := header

	fmt.Println("Private JWT Token String: " + tknStr)

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Send the new json token as a response
	sendResponse(w, tokenString)

}

//CURRENTLY UNUSED METHOD - WILL BE USED IF WE WANT TO IMPLEMENT NORMAL SIGN IN WITH USER AND PASS
func unused_signin(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Get the existing entry present in the database for the given username
	result := Db.QueryRow("select password from userinfo where username=?", creds.Username)
	if err != nil {
		// If there is an issue with the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// We create another instance of `Credentials` to store the credentials we get from the database
	storedCreds := &Credentials{}
	// Store the obtained password in `storedCreds`
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// If the error is of any other type, send a 500 status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		w.WriteHeader(http.StatusUnauthorized)
	}
	// If we reach this point, that means the users password was correct, and that they are authorized
	// The default 200 status is sent
}

//CURRENTLY UNUSED METHOD - WILL BE USED IF WE WANT TO IMPLEMENT NORMAL SIGN UP WITH USER AND PASS
func unused_signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I hit the signup function")
	// Parse and decode the request body into a new `Credentials` instance
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad parameters")
		return
	}
	fmt.Println("I decoded the json")

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		// If there is something wrong with hashing of the password, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("I fail at hashing password")

		return
	}

	// Next, insert the username, along with the hashed password into the database
	if _, err = Db.Query("insert into userinfo (username, password) values(?, ?)", creds.Username, string(hashedPassword)); err != nil {
		// If there is any issue with inserting into the database, return a 500 error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("I fail at inserting val into database")
		fmt.Println(err.Error())
		return
	}
	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}
