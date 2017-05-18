package main

import (

	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
	
)

var logger = shim.NewLogger("PLVChaincode")

//=======================================================================================================================
// Index name
//=======================================================================================================================

const UsersIndexName   =   "users"
const ImagesIndexName  =   "images"

var indexNames = []string{
	UsersIndexName,
	ImagesIndexName,
}


//=======================================================================================================================
// Structure definitions 
//=======================================================================================================================
// SimpleChaincode - A blank struct for use with Shim.
//=======================================================================================================================

type SampleChaincode struct {
}

//=======================================================================================================================
// Image - Defines the structure for an image object.
//=======================================================================================================================

type Image struct {

    ID 		    	string      `json:"id"`
	Name        	string      `json:"name"`
	Author			string		`json:"author"`
	URL				string      `json:"url"`
	User        	string      `json:"user"`
	MD5Hash      	string      `json:"md5-hash"`
	Remarks     	string      `json:"remarks"`
	PurchaseDate	string      `json:"purchase-date"`
	Status          int         `json:"status"`
	
} 

//=======================================================================================================================
// User - participant type could be Empoloyee or Marketing		   
//=======================================================================================================================

type User struct {

	Username        string      `json:"username"`
	Password 		string 		`json:"password"`
	PType           string      `json:"participant-type"`

}

//=======================================================================================================================
// Authentication result 		   
//=======================================================================================================================

type UserAuthenticationResult struct {

	User        	User
	Authenticated 	bool
	
}

//=======================================================================================================================
// Users 		   
//=======================================================================================================================

type Users struct {

	Users 	[]User	 `json:"users"`
	
}

//=======================================================================================================================
// Images		   
//=======================================================================================================================


type Images struct {

	Images 	[]Image	 `json:"images"`
	
}

//=======================================================================================================================
//  Query Functions
//=======================================================================================================================
//  Get index - Get index of images or users
//=======================================================================================================================

func GetIndex(stub shim.ChaincodeStubInterface, indexName string) ([]string, error) {

	indexAsBytes, err := stub.GetState(indexName)
	
	if err != nil {
	
		return nil, errors.New("Failed to get " + indexName)
		
	}

	var index []string
	
	err = json.Unmarshal(indexAsBytes, &index)
	
	if err != nil {
	
		return nil, errors.New("Error unmarshalling index '" + indexName + "': " + err.Error())
		
	}

	return index, nil
	
}


//=======================================================================================================================
//  Add ID to the index
//=======================================================================================================================

func AddIDToIndex(stub shim.ChaincodeStubInterface, indexName string, id string) ([]byte, error) {

	result, err := DoesIDExist(stub, id, indexName)
	
	if result == true {
		
		return nil, errors.New("ID already exists")
	
	}

	index, err := GetIndex(stub, indexName)
	
	if err != nil {
	
		return nil, err
		
	}

	index = append(index, id)

	jsonAsBytes, err := json.Marshal(index)
	
	if err != nil {
	
		return nil, errors.New("Error marshalling index '" + indexName + "': " + err.Error())
		
	}

	err = stub.PutState(indexName, jsonAsBytes)
	
	if err != nil {
	
		return nil, errors.New("Error storing new " + indexName + " into ledger")
		
	}

	return []byte(id), nil
	
}

//=======================================================================================================================
//  Store the object in the ledger 
//=======================================================================================================================

func Store(stub shim.ChaincodeStubInterface, objectID string, indexName string, object []byte) error {

	ID, err := AddIDToIndex(stub, indexName, objectID)
	
	if err != nil {
	
		return errors.New("Writing ID to index: " + indexName + "Reason: " + err.Error())
		
	}

	fmt.Println("adding: ", string(object))

	err = stub.PutState(string(ID), object)
	
	if err != nil {
	
		return errors.New("Putstate error: " + err.Error())
		
	}

	return nil
}

//=======================================================================================================================
//  Add user
//=======================================================================================================================

func addUser(stub shim.ChaincodeStubInterface, index string, userJSONObject string) error {

	id, err := AddIDToIndex(stub, UsersIndexName, index)
	
	if err != nil {
	
		return errors.New("Error creating new id for user " + index)
		
	}

	err = stub.PutState(string(id), []byte(userJSONObject))
	
	if err != nil {
	
		return errors.New("Error putting user data on ledger")
		
	}

	return nil
}

//=======================================================================================================================
//  Demand Image function 
//=======================================================================================================================

func DemandImage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {


	if len(args) != 1 {
	
        logger.Debug("Invalid number of args")
        return nil, errors.New("Expected at least one argument for demanding new image")
		
    }
	
	imageAsJSON := args[0]
	
	var image Image
	
	if err := json.Unmarshal([]byte(imageAsJSON), &image); err != nil {
	
		return nil, errors.New("Error while unmarshalling image, reason: " + err.Error())
		
	}
	
	imageAsBytes, err := json.Marshal(image)
	
	if err != nil {
	
		return nil, errors.New("Error marshalling image, reason: " + err.Error())
		
	}
	
	Store(stub, image.ID, ImagesIndexName, imageAsBytes)
	
	return nil, nil

}

//=======================================================================================================================
//  Deliver Image function 
//=======================================================================================================================

func DeliverImage(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){

    if len(args) < 4 {
	
        logger.Debug("Invalid number of args")
        return nil, errors.New("Expected at least four arguments for delivering new image")
		
    }
 
    var imageId   		=  args[0]
	var Name            =  args[1]
    var MD5Hash   		=  args[2]
	var PurchaseDate    =  args[3]
	
	imageBytes, err := stub.GetState(imageId)
	
	if err != nil {
	
		logger.Error("Could not fetch the image from ledger", err)
		return nil, err
		
	}
	
	var image Image 
	err = json.Unmarshal(imageBytes, &image)
	image.MD5Hash = MD5Hash
	image.PurchaseDate = PurchaseDate
	image.Name = Name
	image.Status = 2
	
	imageBytes, err = json.Marshal(&image)
	
	if err != nil {
	
		logger.Error("Could not marshal image post update", err)
		return nil, err
		
	}
	
	err = stub.PutState(imageId, imageBytes)
	
	if err != nil {
	
		logger.Error("Could not save image post update", err)
		return nil, err
		
	}
 
    fmt.Println("The image is successfully delivered")
	
    return nil, nil

}

//=======================================================================================================================
//  Query Functions
//=======================================================================================================================
//  Get User
//=======================================================================================================================

func GetUser(stub shim.ChaincodeStubInterface, username string) (User, error) {

	userAsBytes, err := stub.GetState(username)
	
	if err != nil {
	
		return User{}, errors.New("Could not retrieve information for this user")
		
	}

	var user User
	
	if err = json.Unmarshal(userAsBytes, &user); err != nil {
	
		return User{}, errors.New("Cannot get user, reason: " + err.Error())
		
	}

	return user, nil
}

//=======================================================================================================================
//   Authenticate User
//=======================================================================================================================

func AuthenticateAsUser(stub shim.ChaincodeStubInterface, user User, password string) (UserAuthenticationResult) {

	if user == (User{}) {
	
		fmt.Println("User not found")
		
		return UserAuthenticationResult{
		
			User: user,
			Authenticated: false,
			
		}
		
	}

	if user.Password != password {
	
		fmt.Println("Password does not match")
		
		return UserAuthenticationResult{
		
			User: user,
			Authenticated: false,
			
		}
		
	}

	return UserAuthenticationResult{
	
		User: user,
		
		Authenticated: true,
		
	}
	
}

//=======================================================================================================================
//  Get Image
//=======================================================================================================================

func getImage(stub shim.ChaincodeStubInterface, imageID string) ([]byte, error) {
    fmt.Println("Get image zone")
 
    if imageID == "" {
        fmt.Println("Invalid number of arguments")
        return nil, errors.New("Missing image ID")
    } 
    bytes, err := stub.GetState(imageID)
    if err != nil {
        fmt.Println("Could not fetch an image with the demand id "+imageID+" from ledger", err)
        return nil, err
    }
	
    return bytes, nil
}

//=======================================================================================================================
//  Get All Images By User Function 
//=======================================================================================================================

func GetImagesByUser(stub shim.ChaincodeStubInterface, User string) ([]byte, error) {

	imagesIndex, err := GetIndex(stub, ImagesIndexName)
	
	if err != nil {
	
		return nil, errors.New("Unable to retrieve imagesIndex, reason: " + err.Error())
		
	}

	// imageIDs := []string{}
	
	var images []Image
	
	for _, imageID := range imagesIndex {
	
		imageAsBytes, err := stub.GetState(imageID)
		
		if err != nil {
		
			return nil, errors.New("Could not retrieve image for ID " + imageID + " reason: " + err.Error())
			
		}

		var image Image
		
		err = json.Unmarshal(imageAsBytes, &image)
		
		if err != nil {
		
			return nil, errors.New("Error while unmarshalling imageAsBytes, reason: " + err.Error())
			
		}
		
		

		if image.User == User {
		
			// imageIDs = append(imageIDs, strconv.Itoa(image.ID))
			
			images = append(images, image);
			
			
			
			
		}
		
	}

	// return imageIDs, nil
	
	return json.Marshal(Images {Images: images})
	
}

//=======================================================================================================================
//  Get All users as objects
//=======================================================================================================================

func GetAllUsers(stub shim.ChaincodeStubInterface) ([]User, error) {

	usersIndex, err := GetIndex(stub, UsersIndexName)
	
	if err != nil {
	
		return []User{}, errors.New("Could not retrieve userIndex, reason: " + err.Error())
		
	}

	var users []User
	
	for _, userID := range usersIndex {
	
		userAsBytes, err := stub.GetState(userID)
		
		if err != nil {
		
			return []User{}, errors.New("Could not retrieve user with ID: " + userID + ", reason: " + err.Error())
			
		}

		var user User
		
		err = json.Unmarshal(userAsBytes, &user)
		
		if err != nil {
		
			return []User{}, errors.New("Error while unmarshalling user, reason: " + err.Error())
			
		}
		
		user.Username = userID

		users = append(users, user)
		
	}

	return users, nil
}

//=======================================================================================================================
//  Get all users as bytes
//=======================================================================================================================

func GetUsers(stub shim.ChaincodeStubInterface) ([]byte, error) {

	users, err := GetAllUsers(stub)
	
	if err != nil {
	
		return nil, err
		
	}

	return json.Marshal( Users {Users: users})
	
}

//=======================================================================================================================
//  Get images as objects
//=======================================================================================================================

func GetAllImages(stub shim.ChaincodeStubInterface) ([]Image, error) {

	imagesIndex, err := GetIndex(stub, ImagesIndexName)
	
	if err != nil {
	
		return []Image{}, errors.New("Could not retrieve ImagesIndex, reason: " + err.Error())
		
	}

	var images []Image
	
	for _, imageID := range imagesIndex {
	
		imageAsBytes, err := stub.GetState(imageID)
		
		if err != nil {
		
			return []Image{}, errors.New("Could not retrieve image with ID: " + imageID + ", reason: " + err.Error())
			
		}

		var image Image
		
		err = json.Unmarshal(imageAsBytes, &image)
		
		if err != nil {
		
			return []Image{}, errors.New("Error while unmarshalling image, reason: " + err.Error())
			
		}

		images = append(images, image)
		
	}

	return images, nil
}

//=======================================================================================================================
//  Get all images as bytes 
//=======================================================================================================================

func GetImages(stub shim.ChaincodeStubInterface) ([]byte, error) {

	images, err := GetAllImages(stub)
	
	if err != nil {
	
		return nil, err
		
	}

	return json.Marshal( Images {Images: images})
	
}

//=======================================================================================================================
//  Check if ID already exists  
//=======================================================================================================================

func DoesIDExist(stub shim.ChaincodeStubInterface, id string, indexName string) (bool, error) {
	index, err := GetIndex(stub, indexName)
	if err != nil {
		return false, err
	}

	for _, indexElement := range index {
		if indexElement == id {
			return true, nil
		}
	}

	return false, nil
}



//#######################################################################################################################
//#																														#
//#                  				Main functions (Init, Invoke and Query) 											#
//#																														#
//#######################################################################################################################

//=======================================================================================================================
//   Init function - Called when the user deploys the chaincode.
//=======================================================================================================================

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string)([]byte, error) {

	for _, indexName := range indexNames {
		var emptyIndex []string

		empty, err := json.Marshal(emptyIndex)
		if err != nil {
			return nil, errors.New("Error marshalling")
		}

		err = stub.PutState(indexName, empty);
		if err != nil {
			return nil, errors.New("Error deleting index")
		}

		logger.Infof("Delete with success from ledger: " + indexName)
	}
	return nil, nil
	
}

//=======================================================================================================================
//  Invoke function
//=======================================================================================================================

func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "addUser" {
	
		// args[0] = new User ID (username)
		// args[1] = new User Data (password, ptype)
		return nil, addUser(stub, args[0], args[1])
	
	} else if  function == "DemandImage" {
	
		return DemandImage(stub, args)
		
	} else if  function == "DeliverImage" {
	
		return DeliverImage(stub, args)
		
	}
	
	return nil, nil
}

//=======================================================================================================================
//  Query Function
//=======================================================================================================================

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
									      
		if function == "AuthenticateAsUser" {
		
			user, err := GetUser(stub, args[0])
			if err != nil {
				logger.Infof("User with id %v not found.", args[0])
			}
		
		    return json.Marshal(AuthenticateAsUser(stub, user, args[1]))
			
		} else if function == "getUsers" {
		
			return GetUsers(stub)
			 
		} else if function == "GetImagesByUser" {
		
			// args[0] : username
			return GetImagesByUser(stub , args[0])
			
		} else if function == "getImage" {
		
			// args[0] : imageID
			return getImage(stub, args[0]);
			
		} else if function == "GetImages" {
		
			return GetImages(stub)
			
		}
		
	return nil, nil
}

//=======================================================================================================================
//  Main Function
//=======================================================================================================================

func main() {

	lld, _ := shim.LogLevel("DEBUG")
	fmt.Println(lld)

	logger.SetLevel(lld)
	fmt.Println(logger.IsEnabledFor(lld))
	
	
    err := shim.Start(new(SampleChaincode))
    if err != nil {
        fmt.Println("Could not start SampleChaincode")
    } else {
        fmt.Println("SampleChaincode successfully started")
    }
 
}
