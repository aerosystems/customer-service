package services

import (
	"errors"
	"fmt"
	"github.com/aerosystems/user-service/internal/models"
	"github.com/aerosystems/user-service/internal/validators"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type CodeService interface {
	CheckCode(code string) error
}

type CodeServiceImpl struct {
	codeRepo models.CodeRepository
}

func NewCodeServiceImpl(codeRepo models.CodeRepository) *CodeServiceImpl {
	return &CodeServiceImpl{
		codeRepo: codeRepo,
	}
}

func (cs *CodeServiceImpl) CheckCode(code string) error {
	if err := validators.ValidateCode(code); err != nil {
		return err
	}
	codeObj, err := cs.codeRepo.GetByCode(code)
	if err != nil {
		return errors.New("could not get data by code")
	}
	if codeObj == nil {
		return errors.New("code does not exist")
	}
	if codeObj.ExpireAt.Before(time.Now()) {
		return errors.New("code is expired")
	}
	if codeObj.IsUsed {
		return errors.New("code is already used")
	}
	return nil
}

func NewCode(UserId uint, Action string, Data string) *models.Code {
	codeExpMinutes, _ := strconv.Atoi(os.Getenv("CODE_EXP_MINUTES"))

	code := models.Code{
		Code:      genCode(),
		UserId:    UserId,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(time.Minute * time.Duration(codeExpMinutes)),
		Action:    Action,
		Data:      Data,
		IsUsed:    false,
	}
	return &code
}

func genCode() string {
	rand.Seed(time.Now().UnixNano())
	var availableNumbers [3]int
	for i := 0; i < 3; i++ {
		availableNumbers[i] = rand.Intn(9)
	}
	var code string
	for i := 0; i < 6; i++ {
		randNum := availableNumbers[rand.Intn(len(availableNumbers))]

		code = fmt.Sprintf("%s%d", code, randNum)
	}
	return code
}
