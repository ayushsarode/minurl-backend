package utils


import 
 (
	    "github.com/gin-gonic/gin"

 )


 func UploadProfilePic(c *gin.Context) (string, error) {
	file,err := c.FormFile("profile_pic")

	if err != nil {
		return "https://avatar.iran.liara.run/public/boy?username=Ash", nil
	}

	filePath := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
 }

