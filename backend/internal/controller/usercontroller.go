package controllers

import (
    "lehrium-backend/internal/models"
    "net/http"
    "github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
    var user models.User
    if err := context.ShouldBindJSON(&user); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        context.Abort()
        return
    }
    if err := user.HashPassword(user.Password); err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        context.Abort()
        return
    }
    // db command go here
    context.JSON(http.StatusCreated, gin.H{"email": user.Email, "untisName": user.UntisName, "password": user.Password })
}
