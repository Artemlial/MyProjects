package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"os"
	"io"
	"log"
)

func encrypt(keStr string, plain []byte) (string,error){
	key,err:=hex.DecodeString(keStr)
	if err!=nil{
		return "",err
	}

	block,err:=aes.NewCipher(key)
	if err!=nil{
		return "",err
	}
	ciphertext:= make([]byte,aes.BlockSize+len(plain))
	iv:=ciphertext[:aes.BlockSize]
	if _,err:=io.ReadFull(rand.Reader,iv);err!=nil{
		return "",err
	}
	stream:=cipher.NewCFBEncrypter(block,iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:],plain)
	return base64.URLEncoding.EncodeToString(ciphertext),nil
}

func main(){
	plain,err := os.ReadFile("")
	if err!=nil{
		panic(err)
	}
	key:=make([]byte,32)
	if _,err:=rand.Read(key);err!=nil{
		log.Fatalf("cannot zahuyachit' enc key : %s\n",err.Error())
	}

	keyStr:=hex.EncodeToString(key)

	cryptotext,err:= encrypt(keyStr,plain)
	if err!=nil{
		log.Fatalf("encryption error: %s",err.Error())
	}
	if err =os.WriteFile("",[]byte(cryptotext),0777);err!=nil{
		log.Printf("could not write enc data: %s\n",err.Error())
	}

}
