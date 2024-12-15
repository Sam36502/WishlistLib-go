package wishlistlib

import (
	"testing"
)

const (
	API_BASE_URL = "https://wishlist-api-wishlist-backend.pearcenet.ch"
	API_PORT     = 0

	TEST_USER_NAME  = "Jim Test"
	TEST_USER_EMAIL = "jim.test@example.com"
	TEST_USER_PASS  = "beowulf"
	TEST_ITEM_NAME  = "Test Item"
	TEST_ITEM_DESC  = "An item to test the API"
	TEST_ITEM_PRICE = 10
)

func TestConnection(t *testing.T) {

	wc := DefaultWishClient(API_BASE_URL, API_PORT)
	wc.Port = 0

	t.Log("Testing the connection...")
	_, err := wc.executeRequest("GET", "/status", nil, nil, false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Connected successfully")
}

func TestCRUD(t *testing.T) {

	var err error
	var user User
	var users []User
	var item Item
	var items []Item

	wc := DefaultWishClient(API_BASE_URL, API_PORT)

	t.Log("Getting all users...")
	users, err = wc.GetAllUsers()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	logUserList(t, users)

	t.Log("Creating new user...")
	user, err = wc.AddNewUser(User{
		Name:  TEST_USER_NAME,
		Email: TEST_USER_EMAIL,
	}, TEST_USER_PASS)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Getting new guy by email...")
	user, err = wc.GetUserByEmail(TEST_USER_EMAIL)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Authenticating...")
	err = wc.Authenticate(user.Email, TEST_USER_PASS)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Log("Getting all of Guy's items...")
	items, err = wc.GetAllItemsOfUser(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	logItemList(t, items)

	t.Log("Add New Item...")
	item, err = wc.AddItemOfUser(Item{
		Name:        TEST_ITEM_NAME,
		Description: TEST_ITEM_DESC,
		Price:       TEST_ITEM_PRICE,
	}, user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("Last Added Item:", item)

	items, _ = wc.GetAllItemsOfUser(user)
	logItemList(t, items)

	t.Log("Reserving new item...")
	err = wc.ReserveItemOfUser(item, user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	items, _ = wc.GetAllItemsOfUser(user)
	logItemList(t, items)

	item, _ = wc.GetItemByID(item.ItemID)
	t.Log("Reserved By: ", item.ReservedByUser.Name)

	t.Log("Un-Reserving new item...")
	err = wc.UnreserveItemOfUser(item, user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	items, _ = wc.GetAllItemsOfUser(user)
	logItemList(t, items)

	t.Log("Deleting the new item...")
	wc.DeleteItemOfUser(item, user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	items, _ = wc.GetAllItemsOfUser(user)
	logItemList(t, items)

	t.Log("Changing user's name...")
	err = wc.ChangeUser(user, "Fred Test", "", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	users, _ = wc.GetAllUsers()
	logUserList(t, users)

	t.Log("Deleting user...")
	err = wc.DeleteUser(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	users, _ = wc.GetAllUsers()
	logUserList(t, users)

}

// util funcs

func logUserList(t *testing.T, us []User) {
	t.Log("  All Users:")
	for _, u := range us {
		t.Log("   - " + u.String())
	}
}

func logItemList(t *testing.T, is []Item) {
	t.Log("  All Items:")
	for _, i := range is {
		t.Log("   - " + i.String())
	}
}
