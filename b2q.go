package main

import (
	"fmt"
	"github.com/zeebo/bencode"
	"io/ioutil"
	"os"
	//"reflect"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"
	"sync"
	//"github.com/davecgh/go-spew/spew"
	"log"
	"bufio"
	//"unicode/utf8"
	"strconv"
)

func decodetorrentfile(path string) map[string]interface{} {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var torrent map[string]interface{}
	if err := bencode.DecodeBytes([]byte(dat), &torrent); err != nil {
		log.Fatal(err)
	}
	return torrent
}

func encodetorrentfile(path string, newstructure map[string]interface{}) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.Create(path)
	}

	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bufferedWriter := bufio.NewWriter(file)
	enc := bencode.NewEncoder(bufferedWriter)
	if err := enc.Encode(newstructure); err != nil {
		log.Fatal(err)
	}
	bufferedWriter.Flush()
	return true
}

func gethash(info interface{}) string {
	torinfo, _ := bencode.EncodeString(info.(map[string]interface{}))
	h := sha1.New()
	io.WriteString(h, torinfo)
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

func piecesconvert(s []byte ) (newpieces []byte) {

	for _, c := range s {
		var binString string
		binString = fmt.Sprintf("%s%b",binString, c)
		for _, d := range binString {
			chr, _ := strconv.Atoi(string(d))
			newpieces = append(newpieces, byte(chr))
		}
	}
	return
}

func logic(key string, value interface{}, bitdir *string, wg *sync.WaitGroup) {
	defer wg.Done()
	newstructure := map[string]interface{}{"active_time": 0, "added_time": 0, "announce_to_dht": 0,
		"announce_to_lsd": 0, "announce_to_trackers": 0, "auto_managed": 0,
		"banned_peers": new(string), "banned_peers6": new(string), "blocks per piece": 0,
		"completed_time": 0, "download_rate_limit": 0, "file sizes": new([][]int),
		"file-format": "libtorrent resume file", "file-version": 0, "file_priority": new([]int), "finished_time": 0,
		"info-hash": new([]byte), "last_seen_complete": 0, "libtorrent-version": "1.1.5.0",
		"max_connections": 0, "max_uploads": 0, "num_complete": 0, "num_downloaded": 0,
		"num_incomplete": 0, "paused": 0, "peers": new(string), "peers6": new(string),
		"pieces": new([]byte), "qBt-category": new(string), "qBt-hasRootFolder": 0, "qBt-name": new(string),
		"qBt-queuePosition": 0, "qBt-ratioLimit": 0, "qBt-savePath": new(string),
		"qBt-seedStatus": 1, "qBt-seedingTimeLimit": -2, "qBt-tags": new([]string),
		"qBt-tempPathDisabled": 0, "save_path": new(string), "seed_mode": 0, "seeding_time": 0,
		"sequential_download": 0, "super_seeding": 0, "total_downloaded": 0,
		"total_uploadedv": 0, "trackers": new([][]string), "upload_rate_limit": 0,
	}
	local := value.(map[string]interface{})
	uncastedlabel := local["labels"].([]interface{})
	newlabels := make([]string, len(uncastedlabel), len(uncastedlabel)+1)
	for num, value := range uncastedlabel {
		if value == nil {
			value = "Empty"
		}
		newlabels[num] = value.(string)
	}
	torrentfile := decodetorrentfile(*bitdir + key)
	var files []string
	if val, ok := torrentfile["info"].(map[string]interface{})["files"].([]interface{}); ok {
		for _, i := range val {
			pathslice := i.(map[string]interface{})["path"]
			var newpath []string
			for _, path := range pathslice.([]interface{}) {
				newpath = append(newpath, path.(string))
			}
			files = append(files, strings.Join(newpath, "/"))
		}
	}
	newstructure["added_time"] = local["added_on"]
	newstructure["completed_time"] = local["completed_on"]
	newstructure["info-hash"] = local["info"]
	newstructure["qBt-tags"] = local["labels"]
	newstructure["blocks per piece"] = torrentfile["info"].(map[string]interface{})["piece length"].(int64) / local["blocksize"].(int64)
	newstructure["pieces"] = piecesconvert([]byte(local["have"].(string)))
	encodetorrentfile("F:/test.fastdecode", newstructure)
}

func main() {
	var wg sync.WaitGroup
	bitdir := "C:/Users/rumanzo/AppData/Roaming/BitTorrent/"
	bitfile := bitdir + "resume.dat"

	torrent := decodetorrentfile(bitfile)

	for key, value := range torrent {
		if key != ".fileguard" && key != "rec" && key == "annabelle-lane-3840x1920.mp4.torrent" {
			wg.Add(1)
			go logic(key, value, &bitdir, &wg)
		}
	}
	wg.Wait()
}
