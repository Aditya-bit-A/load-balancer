package models

import (
	"hash/fnv"
	"os"
	"strconv"

	"github.com/google/uuid"
)

func H(key string) int {
	i := int(HashStringToInt(key))
	K, _ := strconv.Atoi(os.Getenv("K"))
	hash := (i*i + 2*i + 17) % K
	return hash
}

func SH(key string, j int) int {
	i := int(HashStringToInt(key))
	K, _ := strconv.Atoi(os.Getenv("K"))
	serverHash := (i*i + j*j + 2*j + 25) % K
	return serverHash
}

// HashStringToInt takes a string and returns its hashed integer value
func HashStringToInt(s string) uint32 {
	hasher := fnv.New32a()  // 32-bit FNV-1a hash
	hasher.Write([]byte(s)) // write the string bytes to the hasher
	return hasher.Sum32()   // get the resulting hash as a uint32
}
func GenerateRequestId() string {
	return uuid.New().String()
}
