package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/big"
	"os"
	"sort"
)

var (
	blockHash = "2af1a6baa9153906b478d352412479f0c8a611fe896499bdaf0905a15cf32bca"
)

//User 用户结构
type User struct {
	UserID     string
	InviteCode string
	Score      *big.Int
}

type userSlice []*User

func (s userSlice) Len() int           { return len(s) }
func (s userSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s userSlice) Less(i, j int) bool { return s[i].Score.Cmp(s[j].Score) > 0 }

//scoreUser 基于种子seed和用户邀请码进行评分
func scoreUser(users userSlice) {
	seed, err := hex.DecodeString(blockHash)
	if err != nil {
		panic("decode blockhash err")
	}

	base := new(big.Int).SetInt64(10 * 1E8)
	for _, user := range users {
		hasher := sha256.New()
		hasher.Write(seed)
		hasher.Write([]byte(user.InviteCode))
		hashBytes := hasher.Sum(nil)

		hexSha1 := hex.EncodeToString(hashBytes)

		intBase16, success := new(big.Int).SetString(hexSha1, 16)
		if !success {
			panic("Failed parsing big Int from hex")
		}

		user.Score = new(big.Int).Mod(intBase16, base)
	}

	sort.Stable(users)
}

func readData(filename string) userSlice {
	var users userSlice

	file, err := os.Open(filename)

	if err != nil {
		panic("open file failed")
	}
	defer file.Close()

	br := bufio.NewReader(file)
	for {
		data, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		user := &User{}
		err = json.Unmarshal(data, user)
		if err != nil {
			panic("read data failed")
		}

		users = append(users, user)
	}

	return users
}

func writeData(filename string, users userSlice) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		panic("open write file failed")
	}
	defer file.Close()

	for _, user := range users {
		data, err := json.Marshal(user)
		if err != nil {
			panic("Marshal data failed")
		}

		_, err = io.WriteString(file, string(data)+"\n")
		if err != nil {
			panic("write data failed")
		}
	}

}

//输入数据格式为:{"UserID":"1", "InviteCode":"abcdefsd"}
func main() {
	users := readData("user.dat")
	scoreUser(users)
	writeData("result.dat", users)
}
