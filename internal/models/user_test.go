package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Update(t *testing.T) {
	name := "test"
	surname := "test"
	location := "test"
	description := "test"
	profilePhoto := "test"
	var user User

	var modifications UserUpdateDto
	modifications.Name = name
	modifications.Surname = surname
	modifications.Location = location
	modifications.Description = description
	modifications.ProfilePhoto = &profilePhoto
	user.Update(modifications)

	assert.Equal(t, name, user.Name)
	assert.Equal(t, surname, user.Surname)
	assert.Equal(t, location, user.Location)
	assert.Equal(t, description, user.Description)
	assert.Equal(t, profilePhoto, *user.ProfilePhoto)
}
