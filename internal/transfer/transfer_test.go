package transfer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/r3labs/diff/v2"
	_ "github.com/r3labs/diff/v2"
	"github.com/rumanzo/bt2qbt/internal/options"
	"github.com/rumanzo/bt2qbt/pkg/qBittorrentStructures"
	"github.com/rumanzo/bt2qbt/pkg/torrentStructures"
	"github.com/rumanzo/bt2qbt/pkg/utorrentStructs"
)

func TestTransferStructure_HandleSavePaths(t *testing.T) {
	type SearchPathCase struct {
		name                 string
		mustFail             bool
		newTransferStructure *TransferStructure
		key                  string
		expected             *TransferStructure
	}

	cases := []SearchPathCase{
		{
			name: "001 Test torrent with windows single nofolder (original) path without replaces",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test_torrent.txt`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: `test_torrent.txt`,
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					QBtContentLayout: "Original",
					Name:             `test_torrent.txt`,
				},
			},
		},
		{
			name: "002 Test torrent with windows single nofolder (original) path with replace",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test_torrent.txt`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent.txt",
					},
				},
				Opts: &options.Opts{PathSeparator: `\`, Replaces: []string{`D:/torrents,E:/newfolder`}},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `E:/newfolder/`,
					SavePath:         `E:\newfolder\`,
					QBtContentLayout: "Original",
					Name:             `test_torrent.txt`,
				},
			},
		},
		{
			name: "003 Test torrent with windows single nofolder (original) path without replaces. NameUTF8",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test_torrent.txt`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						NameUTF8: "test_torrent.txt",
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					QBtContentLayout: "Original",
					Name:             `test_torrent.txt`,
				},
			},
		},
		{
			name: "004 Test torrent with windows single nofolder (original) path without replaces. Renamed File",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\renamed_test_torrent.txt`,
					Targets: [][]interface{}{
						[]interface{}{
							0,
							"renamed_test_torrent.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent.txt",
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					QBtContentLayout: "Original",
					Name:             `test_torrent.txt`,
					MappedFiles:      []string{`renamed_test_torrent.txt`},
				},
			},
		},
		{
			name: "005 Test torrent with windows single nofolder (original) path with replace to linux paths and linux sep. Renamed File",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\renamed_test_torrent.txt`,
					Targets: [][]interface{}{
						[]interface{}{
							0,
							"renamed_test_torrent.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent.txt",
					},
				},
				Opts: &options.Opts{PathSeparator: `/`, Replaces: []string{`D:/torrents,/mnt/d/torrents`}},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `/mnt/d/torrents/`,
					SavePath:         `/mnt/d/torrents/`,
					Name:             `test_torrent.txt`,
					QBtContentLayout: "Original",
					MappedFiles:      []string{`renamed_test_torrent.txt`},
				},
			},
		},
		{
			name: "006 Test torrent with windows folder (original) path without replaces",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test_torrent`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					QBtContentLayout: "Original",
					Name:             `test_torrent`,
				},
			},
		},
		{
			name: "007 Test torrent with windows folder (original) path with replace",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test_torrent`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`, Replaces: []string{`D:/torrents,E:/newfolder`}},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `E:/newfolder/`,
					SavePath:         `E:\newfolder\`,
					QBtContentLayout: "Original",
					Name:             `test_torrent`,
				},
			},
		},
		{
			name: "008 Test torrent with windows folder (original) path without replaces. NameUTF8",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test_torrent`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						NameUTF8: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					QBtContentLayout: "Original",
					Name:             `test_torrent`,
				},
			},
		},
		{
			name: "009 Test torrent with windows folder (original) path without replaces. Renamed File",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\test_torrent`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					QBtContentLayout: "Original",
					Name:             `test_torrent`,
					MappedFiles: []string{
						``,
						``,
						`test_torrent\renamed_test_torrent.txt`,
					},
				},
			},
		},
		{
			name: "010 Test torrent with windows folder (original) path with replace to linux paths and linux sep. Renamed File",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\test_torrent`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `/`, Replaces: []string{`D:/torrents,/mnt/d/torrents`}},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `/mnt/d/torrents/`,
					SavePath:         `/mnt/d/torrents/`,
					QBtContentLayout: "Original",
					Name:             `test_torrent`,
					MappedFiles: []string{
						``,
						``,
						`test_torrent/renamed_test_torrent.txt`,
					},
				},
			},
		},
		{
			name: "011 Test torrent with windows folder (NoSubfolder) path without replaces",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath: `D:/torrents/test`,
					SavePath:    `D:\torrents\test`,
					Name:        `test_torrent`,
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`file0.txt`,
					},
					QBtContentLayout: "NoSubfolder",
				},
			},
		},
		{
			name: "012 Test torrent with windows folder (NoSubfolder) path with replace",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`, Replaces: []string{`D:/torrents,E:/newfolder`}},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath: `E:/newfolder/test`,
					SavePath:    `E:\newfolder\test`,
					Name:        `test_torrent`,
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`file0.txt`,
					},
					QBtContentLayout: "NoSubfolder",
				},
			},
		},
		{
			name: "013 Test torrent with windows folder (NoSubfolder) path without replaces. NameUTF8",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						NameUTF8: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath: `D:/torrents/test`,
					SavePath:    `D:\torrents\test`,
					Name:        `test_torrent`,
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`file0.txt`,
					},
					QBtContentLayout: "NoSubfolder",
				},
			},
		},
		{
			name: "014 Test torrent with windows folder (NoSubfolder) path without replaces. Renamed File",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\test`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/test`,
					SavePath:         `D:\torrents\test`,
					QBtContentLayout: "NoSubfolder",
					Name:             `test_torrent`,
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`renamed_test_torrent.txt`,
					},
				},
			},
		},
		{
			name: "015 Test torrent with windows folder (NoSubfolder) path with replace to linux paths and linux sep. Renamed File",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\test`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `/`, Replaces: []string{`D:/torrents,/mnt/d/torrents`}},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `/mnt/d/torrents/test`,
					SavePath:         `/mnt/d/torrents/test`,
					Name:             `test_torrent`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`dir1/file1.txt`,
						`dir2/file2.txt`,
						`renamed_test_torrent.txt`,
					},
				},
			},
		},
		{
			name: "016 Test torrent with windows folder (NoSubfolder) path without replaces. TorrentPaths UTF8",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{PathUTF8: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{PathUTF8: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{PathUTF8: []string{"file0.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath: `D:/torrents/test`,
					SavePath:    `D:\torrents\test`,
					Name:        `test_torrent`,
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`file0.txt`,
					},
					QBtContentLayout: "NoSubfolder",
				},
			},
		},
		{
			name:     "017 Test torrent with windows single nofolder (original) path without replaces",
			mustFail: true,
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent.txt`},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test_torrent.txt`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: `test_torrent.txt`,
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torre`,
					Name:             `test_torrent.txt`,
					QBtContentLayout: "Original",
				},
			},
		},
		{
			name: "018 Test torrent with windows folder (original) path without replaces. Moved files with absolute paths",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\test_torrent`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
						[]interface{}{
							int64(3),
							`E:\somedir1\renamed_test_torrent2.txt`,
						},
						[]interface{}{
							int64(4),
							`F:\somedir\somedir4\renamed_test_torrent3.txt`,
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					Name:             `test_torrent`,
					QBtContentLayout: "Original",
					MappedFiles: []string{
						``,
						``,
						`test_torrent\renamed_test_torrent.txt`,
						`E:\somedir1\renamed_test_torrent2.txt`,
						`F:\somedir\somedir4\renamed_test_torrent3.txt`,
					},
				},
			},
		},
		{
			name: "019 Test torrent with windows folder (original) path with replaces. Moved files with absolute paths",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\test_torrent`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
						[]interface{}{
							int64(3),
							`E:\somedir1\renamed_test_torrent2.txt`,
						},
						[]interface{}{
							int64(4),
							`F:\somedir\somedir4\renamed_test_torrent3.txt`,
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `/`, Replaces: []string{`D:/torrents,/mnt/d/torrents`, `E:,/mnt/e`, `F:/,/mnt/f/`}},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `/mnt/d/torrents/`,
					SavePath:         `/mnt/d/torrents/`,
					Name:             `test_torrent`,
					QBtContentLayout: "Original",
					MappedFiles: []string{
						``,
						``,
						`test_torrent/renamed_test_torrent.txt`,
						`/mnt/e/somedir1/renamed_test_torrent2.txt`,
						`/mnt/f/somedir/somedir4/renamed_test_torrent3.txt`,
					},
				},
			},
		},
		{
			name: "020 Test torrent with windows folder (NoSubfolder) path without replaces. Moved files with absolute paths",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\test`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
						[]interface{}{
							int64(3),
							`E:\somedir1\renamed_test_torrent2.txt`,
						},
						[]interface{}{
							int64(4),
							`F:\somedir\somedir4\renamed_test_torrent3.txt`,
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/test`,
					SavePath:         `D:\torrents\test`,
					Name:             `test_torrent`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`renamed_test_torrent.txt`,
						`E:\somedir1\renamed_test_torrent2.txt`,
						`F:\somedir\somedir4\renamed_test_torrent3.txt`,
					},
				},
			},
		},
		{
			name: "021 Test torrent with windows share folder (Original) path without replaces. Moved files with absolute paths. Windows share",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `\\torrents\test_torrent`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
						[]interface{}{
							int64(3),
							`E:\somedir1\renamed_test_torrent2.txt`,
						},
						[]interface{}{
							int64(4),
							`\\somedir\somedir4\renamed_test_torrent3.txt`,
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `//torrents/`,
					SavePath:         `\\torrents\`,
					Name:             `test_torrent`,
					QBtContentLayout: "Original",
					MappedFiles: []string{
						``,
						``,
						`test_torrent\renamed_test_torrent.txt`,
						`E:\somedir1\renamed_test_torrent2.txt`,
						`\\somedir\somedir4\renamed_test_torrent3.txt`,
					},
				},
			},
		},
		{
			name: "022 Test torrent with windows share folder (NoSubfolder) path without replaces. Moved files with absolute paths",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `\\torrents\test`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
						[]interface{}{
							int64(3),
							`E:\somedir1\renamed_test_torrent2.txt`,
						},
						[]interface{}{
							int64(4),
							`\\somedir\somedir4\renamed_test_torrent3.txt`,
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `//torrents/test`,
					SavePath:         `\\torrents\test`,
					Name:             `test_torrent`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`renamed_test_torrent.txt`,
						`E:\somedir1\renamed_test_torrent2.txt`,
						`\\somedir\somedir4\renamed_test_torrent3.txt`,
					},
				},
			},
		},
		{
			name: "023 Test torrent with windows folder (Original) path without replaces. Absolute paths. Windows share Replace",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Name: `test_torrent`},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `\\torrents\test_torrent`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
						[]interface{}{
							int64(3),
							`E:\somedir1\renamed_test_torrent2.txt`,
						},
						[]interface{}{
							int64(4),
							`\\somedir\somedir4\renamed_test_torrent3.txt`,
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `/`, Replaces: []string{`D:/torrents,/mnt/d/torrents`, `E:,/mnt/e`, `//somedir,/mnt/share`}},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `//torrents/`,
					SavePath:         `//torrents/`,
					Name:             `test_torrent`,
					QBtContentLayout: "Original",
					MappedFiles: []string{
						``,
						``,
						`test_torrent/renamed_test_torrent.txt`,
						`/mnt/e/somedir1/renamed_test_torrent2.txt`,
						`/mnt/share/somedir4/renamed_test_torrent3.txt`,
					},
				},
			},
		},
		{
			name: "024 Test torrent with windows folder (NoSubfolder) path without replaces. Absolute paths. Windows share Replace",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `\\torrents\test`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"renamed_test_torrent.txt",
						},
						[]interface{}{
							int64(3),
							`E:\somedir1\renamed_test_torrent2.txt`,
						},
						[]interface{}{
							int64(4),
							`\\somedir\somedir4\renamed_test_torrent3.txt`,
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `/`, Replaces: []string{`D:/torrents,/mnt/d/torrents`, `E:,/mnt/e`, `//somedir,/mnt/share`}},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `//torrents/test`,
					SavePath:         `//torrents/test`,
					Name:             `test_torrent`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`dir1/file1.txt`,
						`dir2/file2.txt`,
						`renamed_test_torrent.txt`,
						`/mnt/e/somedir1/renamed_test_torrent2.txt`,
						`/mnt/share/somedir4/renamed_test_torrent3.txt`,
					},
				},
			},
		},
		{
			name: "025 Test magnet link downloads",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\test`,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{Name: "torrentname"},
				},
				Opts:   &options.Opts{PathSeparator: `\`},
				Magnet: true,
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/test`,
					SavePath:         `D:\torrents\test`,
					Name:             "torrentname",
					QBtContentLayout: "Original",
				},
			},
		},
		{
			name: "026 Test torrent with signle file torrent and savepath in rootdirectory",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\test.txt`,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test.txt",
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/`,
					SavePath:         `D:\`,
					Name:             `test.txt`,
					QBtContentLayout: "Original",
				},
			},
		},
		{
			name: "027 Test torrent with multi file torrent and savepath in rootdirectory",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\test_torrent`,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/`,
					SavePath:         `D:\`,
					Name:             `test_torrent`,
					QBtContentLayout: "Original",
				},
			},
		},
		{
			name: "028 Test torrent with windows folder (original) path without replaces. Emoji utf8 in file and name",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: "D:\\torrents\\test_torrent \xf0\x9f\x86\x95",
					Targets: [][]interface{}{
						[]interface{}{
							int64(0),
							"E:\\somedir1 \xf0\x9f\x86\x95\\\xf0\x9f\x86\x95 renamed_test_torrent2.txt",
						},
						[]interface{}{
							int64(1),
							"\\\\somedir\\somedir4 \xf0\x9f\x86\x95\\\xf0\x9f\x86\x95 renamed_test_torrent3.txt",
						},
						[]interface{}{
							int64(2),
							"renamed \xf0\x9f\x86\x95 file1.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent \xf0\x9f\x86\x95",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "\xf0\x9f\x86\x95 file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"\xf0\x9f\x86\x95", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0 \xf0\x9f\x86\x95.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					Name:             "test_torrent \xf0\x9f\x86\x95",
					QBtContentLayout: "Original",
					MappedFiles: []string{
						"E:\\somedir1 \xf0\x9f\x86\x95\\\xf0\x9f\x86\x95 renamed_test_torrent2.txt",
						"\\\\somedir\\somedir4 \xf0\x9f\x86\x95\\\xf0\x9f\x86\x95 renamed_test_torrent3.txt",
						"test_torrent \xf0\x9f\x86\x95\\renamed \xf0\x9f\x86\x95 file1.txt",
					},
				},
			},
		},
		{
			name: "029 Test torrent with windows folder (NoSubfolder) with renamed files. Emoji utf8 in file and torrent name",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: "D:\\torrents\\renamed test_torrent \xf0\x9f\x86\x95"},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent \xf0\x9f\x86\x95",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "\xf0\x9f\x86\x95 file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"\xf0\x9f\x86\x95", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0 \xf0\x9f\x86\x95.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      "D:/torrents/renamed test_torrent \xf0\x9f\x86\x95",
					SavePath:         "D:\\torrents\\renamed test_torrent \xf0\x9f\x86\x95",
					Name:             "test_torrent \xf0\x9f\x86\x95",
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						"dir1\\\xf0\x9f\x86\x95 file1.txt",
						"\xf0\x9f\x86\x95\\file2.txt",
						"file0 \xf0\x9f\x86\x95.txt",
					},
				},
			},
		},
		{
			name: "030 Test torrent with windows folder (original) path with renamed files. Emoji cesu8 in file and name",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: "D:\\torrents\\test_torrent \xed\xa0\xbc\xed\xb6\x95",
					Targets: [][]interface{}{
						[]interface{}{
							int64(0),
							"E:\\somedir1 \xed\xa0\xbc\xed\xb6\x95\\\xed\xa0\xbc\xed\xb6\x95 renamed_test_torrent2.txt",
						},
						[]interface{}{
							int64(1),
							"\\\\somedir\\somedir4 \xed\xa0\xbc\xed\xb6\x95\\\xed\xa0\xbc\xed\xb6\x95 renamed_test_torrent3.txt",
						},
						[]interface{}{
							int64(2),
							"renamed \xed\xa0\xbc\xed\xb6\x95 file1.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent \xed\xa0\xbc\xed\xb6\x95",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "\xed\xa0\xbc\xed\xb6\x95 file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"\xed\xa0\xbc\xed\xb6\x95", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0 \xed\xa0\xbc\xed\xb6\x95.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      "D:/torrents/test_torrent \xf0\x9f\x86\x95",
					SavePath:         "D:\\torrents\\test_torrent \xf0\x9f\x86\x95",
					Name:             "test_torrent \xf0\x9f\x86\x95",
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						"E:\\somedir1 \xf0\x9f\x86\x95\\\xf0\x9f\x86\x95 renamed_test_torrent2.txt",
						"\\\\somedir\\somedir4 \xf0\x9f\x86\x95\\\xf0\x9f\x86\x95 renamed_test_torrent3.txt",
						"renamed \xf0\x9f\x86\x95 file1.txt",
					},
				},
			},
		},
		{
			name: "031 Test torrent with windows folder (NoSubfolder) with renamed files. Emoji cesu8 in file and torrent name",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: "D:\\torrents\\renamed test_torrent \xed\xa0\xbc\xed\xb6\x95"},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent \xed\xa0\xbc\xed\xb6\x95",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "\xed\xa0\xbc\xed\xb6\x95 file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"\xed\xa0\xbc\xed\xb6\x95", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0 \xed\xa0\xbc\xed\xb6\x95.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      "D:/torrents/renamed test_torrent \xf0\x9f\x86\x95",
					SavePath:         "D:\\torrents\\renamed test_torrent \xf0\x9f\x86\x95",
					Name:             "test_torrent \xf0\x9f\x86\x95",
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						"dir1\\\xf0\x9f\x86\x95 file1.txt",
						"\xf0\x9f\x86\x95\\file2.txt",
						"file0 \xf0\x9f\x86\x95.txt",
					},
				},
			},
		},
		{
			name: "032 Test torrent with windows folder without renamed files. Emoji cesu8 in file and torrent name",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: "D:\\torrents\\test_torrent"},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "\xed\xa0\xbc\xed\xb6\x95 file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"\xed\xa0\xbc\xed\xb6\x95", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0 \xed\xa0\xbc\xed\xb6\x95.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2 \xed\xa0\xbc\xed\xb6\x95.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      "D:/torrents/test_torrent",
					SavePath:         "D:\\torrents\\test_torrent",
					Name:             "test_torrent",
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						"dir1\\\xf0\x9f\x86\x95 file1.txt",
						"\xf0\x9f\x86\x95\\file2.txt",
						"file0 \xf0\x9f\x86\x95.txt",
						`file1.txt`,
						"file2 \xf0\x9f\x86\x95.txt",
					},
				},
			},
		},
		{
			name: "033 Test torrent with windows folder with renamed files. Emoji cesu8 in file and torrent name",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: "D:\\torrents\\test_torrent",
					Targets: [][]interface{}{
						[]interface{}{
							int64(0),
							"E:\\somedir1 \xed\xa0\xbc\xed\xb6\x95\\\xed\xa0\xbc\xed\xb6\x95 renamed_test_torrent2.txt",
						},
						[]interface{}{
							int64(1),
							"\\\\somedir\\somedir4 \xed\xa0\xbc\xed\xb6\x95\\\xed\xa0\xbc\xed\xb6\x95 renamed_test_torrent3.txt",
						},
						[]interface{}{
							int64(3),
							"renamed \xed\xa0\xbc\xed\xb6\x95 file1.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "\xed\xa0\xbc\xed\xb6\x95 file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"\xed\xa0\xbc\xed\xb6\x95", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0 \xed\xa0\xbc\xed\xb6\x95.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2 \xed\xa0\xbc\xed\xb6\x95.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file4.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      "D:/torrents/test_torrent",
					SavePath:         "D:\\torrents\\test_torrent",
					Name:             "test_torrent",
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						"E:\\somedir1 \xf0\x9f\x86\x95\\\xf0\x9f\x86\x95 renamed_test_torrent2.txt",
						"\\\\somedir\\somedir4 \xf0\x9f\x86\x95\\\xf0\x9f\x86\x95 renamed_test_torrent3.txt",
						"file0 \xf0\x9f\x86\x95.txt",
						"renamed \xf0\x9f\x86\x95 file1.txt",
						"file2 \xf0\x9f\x86\x95.txt",
						"file4.txt",
					},
				},
			},
		},
		{
			name: "034 Test torrent with windows single nofolder (original) path without replaces. Cesu8 emoji",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: "D:\\torrents\\test_torrent \xed\xa0\xbc\xed\xb6\x95.txt"},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent \xed\xa0\xbc\xed\xb6\x95.txt",
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath: `D:/torrents/`,
					SavePath:    `D:\torrents\`,
					Name:        "test_torrent \xf0\x9f\x86\x95.txt",
					MappedFiles: []string{
						"test_torrent \xf0\x9f\x86\x95.txt",
					},
					QBtContentLayout: "Original",
				},
			},
		},
		{
			name: "035 Test torrent with windows folder (NoSubfolder) with special symbols.",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\renamed test_torrent`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`#test _ test [01]{1} [6K].jpg`}},
							&torrentStructures.TorrentFile{Path: []string{`testdir1 collection`, `testdir2_`, `1.jpg`}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/renamed test_torrent`,
					SavePath:         `D:\torrents\renamed test_torrent`,
					Name:             `test_torrent`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`#test _ test [01]{1} [6K].jpg`,
						`testdir1 collection\testdir2_\1.jpg`,
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
		},
		{
			name: "035 Test torrent with windows folder (NoSubfolder) with prohibited symbols.",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\renamed test_torrent`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`#test | test [01]{1} [6K].jpg`}},
							&torrentStructures.TorrentFile{Path: []string{`testdir1 collection`, `testdir2?`, `1.jpg`}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/renamed test_torrent`,
					SavePath:         `D:\torrents\renamed test_torrent`,
					Name:             `test_torrent`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`#test _ test [01]{1} [6K].jpg`,
						`testdir1 collection\testdir2_\1.jpg`,
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
		},
		{
			name: "036 Test torrent with windows single nofolder (original) path without replaces. With prohibited symbols",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test_torrent.txt`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: `test|torrent.txt`,
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					Name:             `test_torrent.txt`,
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					QBtContentLayout: "Original",
					MappedFiles: []string{
						`test_torrent.txt`,
					},
				},
			},
		},
		{
			name: "037 Test torrent with windows folder (NoSubfolder) with prohibited symbols. *nix path separator",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\renamed test_torrent`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`#test | test [01]{1} [6K].jpg`}},
							&torrentStructures.TorrentFile{Path: []string{`testdir1 collection`, `testdir2?`, `1.jpg`}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `/`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/renamed test_torrent`,
					SavePath:         `D:/torrents/renamed test_torrent`,
					Name:             `test_torrent`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`#test _ test [01]{1} [6K].jpg`,
						`testdir1 collection/testdir2_/1.jpg`,
					},
				},
			},
		},
		{
			name: "038 Test torrent with windows single nofolder (original) path without replaces. With prohibited symbols. *nix path separator",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Path: `D:\torrents\test_torrent.txt`},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: `test|torrent.txt`,
					},
				},
				Opts: &options.Opts{PathSeparator: `/`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:/torrents/`,
					Name:             `test_torrent.txt`,
					QBtContentLayout: "Original",
					MappedFiles: []string{
						`test_torrent.txt`,
					},
				},
			},
		},
		{
			name: "039 Test torrent with windows single nofolder (original) path with renamed file. v2 torrent",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\testtorrent.txt`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(0),
							"testtorrent.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: `testtorrent.txt`,
						FileTree: map[string]interface{}{
							`testtorrent.txt`: map[string]interface{}{
								``: map[string]interface{}{
									`length`:      int64(100),
									`pieces_root`: []byte{},
								},
							},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					Name:             `testtorrent.txt`,
					QBtContentLayout: "Original",
				},
			},
		},
		{
			name: "040 Test torrent with windows folder (original) path with renamed file. v2 torrent",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\torrents\testtorrent`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(1),
							"testtorrent2.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: `testtorrent`,
						FileTree: map[string]interface{}{
							`testtorrent1.txt`: map[string]interface{}{
								``: map[string]interface{}{
									`length`:      int64(100),
									`pieces_root`: []byte{},
								},
							},
							`testtorrent3.txt`: map[string]interface{}{
								``: map[string]interface{}{
									`length`:      int64(100),
									`pieces_root`: []byte{},
								},
							},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/torrents/`,
					SavePath:         `D:\torrents\`,
					Name:             `testtorrent`,
					QBtContentLayout: "Original",
					MappedFiles: []string{
						``,
						`testtorrent\testtorrent2.txt`,
					},
				},
			},
		},
		{
			name: "041 Test torrent with multi file torrent (NoSubfolder) with prohibited symbols in name",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\test_torrent_test`,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent/test",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/test_torrent_test`,
					SavePath:         `D:\test_torrent_test`,
					Name:             `test_torrent_test`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`file0.txt`,
						`file1.txt`,
						`file2.txt`,
					},
				},
			},
		},
		{
			name: "042 Test torrent with multi file torrent with renames (NoSubFolder) with prohibited symbols in name",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\test_torrent_test`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"file0_file.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent/test",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/test_torrent_test`,
					SavePath:         `D:\test_torrent_test`,
					Name:             `test_torrent_test`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`file0_file.txt`,
						`file1.txt`,
						`file2.txt`,
					},
				},
			},
		},
		{
			name: "043 Test torrent with multi file torrent (NoSubfolder) with space at the end of torrent name",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\test_torrent_`,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent ",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"dir1", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir2", "file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/test_torrent_`,
					SavePath:         `D:\test_torrent_`,
					Name:             `test_torrent_`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`dir1\file1.txt`,
						`dir2\file2.txt`,
						`file0.txt`,
						`file1.txt`,
						`file2.txt`,
					},
				},
			},
		},
		{
			name: "044 Test torrent with multi file torrent (NoSubfolder) with space at the end directory",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\test_torrent_`,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent ",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir1 ", "dir2 ", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir3 ", "file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/test_torrent_`,
					SavePath:         `D:\test_torrent_`,
					Name:             `test_torrent_`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`file0.txt`,
						`file1.txt`,
						`file2.txt`,
						`dir1_\dir2_\file1.txt`,
						`dir3_\file2.txt`,
					},
				},
			},
		},
		{
			name: "045 Test torrent with multi file torrent (NoSubfolder) with space at the end directory with renamed files",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path: `D:\test_torrent_`,
					Targets: [][]interface{}{
						[]interface{}{
							int64(2),
							"file2_renamed.txt",
						},
					},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_torrent ",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"file0.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file2.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir1 ", "dir2 ", "file1.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"dir3 ", "file2.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      `D:/test_torrent_`,
					SavePath:         `D:\test_torrent_`,
					Name:             `test_torrent_`,
					QBtContentLayout: "NoSubfolder",
					MappedFiles: []string{
						`file0.txt`,
						`file1.txt`,
						`file2_renamed.txt`,
						`dir1_\dir2_\file1.txt`,
						`dir3_\file2.txt`,
					},
				},
			},
		},
		{
			name: "046 Test torrent with multi file torrent with transfer to NoSubfolder cesu8 symbols in names",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path:    "D:\\test_slashes_emoji \xed\xa0\xbc\xed\xb6\x95_",
					Caption: "test_slashes_emoji \xed\xa0\xbc\xed\xb6\x95_",
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test_slashes/emoji \xed\xa0\xbc\xed\xb6\x95 ",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"file_with_emoji \xed\xa0\xbc\xed\xb6\x95.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"file_with/slash.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"testdir_with_emoji_and_space \xed\xa0\xbc\xed\xb6\x95 ", "file_with_emoji \xed\xa0\xbc\xed\xb6\x95.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"testdir_with_emoji_and_space \xed\xa0\xbc\xed\xb6\x95 ", "file_with/slash.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"testdir_with_space ", "file_with_emoji \xed\xa0\xbc\xed\xb6\x95.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"testdir_with_space ", "file_with/slash.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      "D:/test_slashes_emoji \xf0\x9f\x86\x95_",
					SavePath:         "D:\\test_slashes_emoji \xf0\x9f\x86\x95_",
					Name:             "test_slashes_emoji \xf0\x9f\x86\x95_",
					QbtName:          "test_slashes_emoji \xf0\x9f\x86\x95_",
					QBtContentLayout: `NoSubfolder`,
					MappedFiles: []string{
						"file_with_emoji \xf0\x9f\x86\x95.txt",
						"file_with_slash.txt",
						"testdir_with_emoji_and_space \xf0\x9f\x86\x95_\\file_with_emoji \xf0\x9f\x86\x95.txt",
						"testdir_with_emoji_and_space \xf0\x9f\x86\x95_\\file_with_slash.txt",
						"testdir_with_space_\\file_with_emoji \xf0\x9f\x86\x95.txt",
						"testdir_with_space_\\file_with_slash.txt",
					},
				},
			},
		},
		{
			name: "047 Test torrent with multi file torrent with transfer to NoSubfolder RLF LRF symbols in torrent name and files",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Path:    "D:\\test files \u200e\u200f",
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Name: "test files \u200e\u200f",
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{"file_with_emoji \u200e\u200f\xed\xa0\xbc\xed\xb6\x95.txt"}},
							&torrentStructures.TorrentFile{Path: []string{"testdir_with_emoji_and_space \u200e\u200f\xed\xa0\xbc\xed\xb6\x95 ", "file_with/slash.txt"}},
						},
					},
				},
				Opts: &options.Opts{PathSeparator: `\`},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					QbtSavePath:      "D:/test files \u200e\u200f",
					SavePath:         "D:\\test files \u200e\u200f",
					Name:             "test files \u200e\u200f",
					QBtContentLayout: `NoSubfolder`,
					MappedFiles: []string{
						"file_with_emoji \u200e\u200f\xf0\x9f\x86\x95.txt",
						"testdir_with_emoji_and_space \u200e\u200f\xf0\x9f\x86\x95_\\file_with_slash.txt",
					},
				},
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.newTransferStructure.Opts != nil {
				replaces := CreateReplaces(testCase.newTransferStructure.Opts.Replaces)
				testCase.newTransferStructure.Replace = replaces
				testCase.expected.Replace = replaces
			}
			testCase.newTransferStructure.Fastresume.Name, _ = testCase.newTransferStructure.TorrentFile.GetNormalizedTorrentName()
			testCase.newTransferStructure.HandleCaption()
			testCase.newTransferStructure.HandleSavePaths()
			equal := reflect.DeepEqual(testCase.expected.Fastresume, testCase.newTransferStructure.Fastresume)
			if !equal && !testCase.mustFail {
				changes, err := diff.Diff(testCase.newTransferStructure.Fastresume, testCase.expected.Fastresume, diff.DiscardComplexOrigin())
				if err != nil {
					t.Error(err.Error())
				}
				t.Fatalf("Unexpected error: opts isn't equal:\nGot: %#v\nExpect %#v\nDiff: %v\n", testCase.newTransferStructure.Fastresume, testCase.expected.Fastresume, spew.Sdump(changes))
			} else if equal && testCase.mustFail {
				t.Fatalf("Unexpected error: structures are equal, but they shouldn't\nGot: %v\n", spew.Sdump(testCase.newTransferStructure.Fastresume))
				fmt.Print("\u200e")
			}
		})
	}
}

func TestTransferStructure_HandlePieces(t *testing.T) {
	type HandlePiecesCase struct {
		name                 string
		mustFail             bool
		newTransferStructure *TransferStructure
		expected             *TransferStructure
	}

	cases := []HandlePiecesCase{
		{
			name: "001 parted",
			newTransferStructure: &TransferStructure{
				NumPieces: 5,
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1, 0, 0},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 5},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 5},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 5},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 5},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 5},
						},
						PieceLength: 5,
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1, 0, 0},
					Pieces: []byte{
						byte(1),
						byte(0),
						byte(1),
						byte(0),
						byte(0),
					},
				},
			},
		},
		{
			name: "002 parted",
			newTransferStructure: &TransferStructure{
				NumPieces: 5,
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 13},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 7},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 5},
						},
						PieceLength: 5, // 25 total
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1},
					Pieces: []byte{
						byte(1),
						byte(1),
						byte(1),
						byte(0),
						byte(1),
					},
				},
			},
		},
		{
			name: "003 parted",
			newTransferStructure: &TransferStructure{
				NumPieces: 5,
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{0, 1, 0},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 9},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 6},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 10},
						},
						PieceLength: 5, // 25 total
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{0, 1, 0},
					Pieces: []byte{
						byte(0),
						byte(1),
						byte(1),
						byte(0),
						byte(0),
					},
				},
			},
		},
		{
			name: "004 parted Mustfail",
			newTransferStructure: &TransferStructure{
				NumPieces: 5,
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1},
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 13},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 7},
							&torrentStructures.TorrentFile{Path: []string{`/`}, Length: 5},
						},
						PieceLength: 5, // 25 total
					},
				},
			},
			mustFail: true,
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1},
					Pieces: []byte{
						byte(0),
						byte(1),
						byte(1),
						byte(0),
						byte(1),
					},
				},
			},
		},
		{
			name: "005 single unfinished",
			newTransferStructure: &TransferStructure{
				NumPieces: 5,
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					Unfinished: new([]interface{}),
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					Unfinished: new([]interface{}),
					Pieces: []byte{
						byte(0),
						byte(0),
						byte(0),
						byte(0),
						byte(0),
					},
				},
			},
		},
		{
			name: "005 single finished",
			newTransferStructure: &TransferStructure{
				NumPieces:  5,
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					Pieces: []byte{
						byte(1),
						byte(1),
						byte(1),
						byte(1),
						byte(1),
					},
				},
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.newTransferStructure.HandlePieces()
			equal := reflect.DeepEqual(testCase.expected.Fastresume, testCase.newTransferStructure.Fastresume)
			if !equal && !testCase.mustFail {
				changes, err := diff.Diff(testCase.newTransferStructure.Fastresume, testCase.expected.Fastresume, diff.DiscardComplexOrigin())
				if err != nil {
					t.Error(err.Error())
				}
				t.Fatalf("Unexpected error: opts isn't equal:\n Got: %#v\n Expect %#v\n Diff: %v\n", testCase.newTransferStructure.Fastresume.Pieces, testCase.expected.Fastresume.Pieces, spew.Sdump(changes))
			} else if equal && testCase.mustFail {
				t.Fatalf("Unexpected error: structures are equal, but they shouldn't\n Got: %v\n", spew.Sdump(testCase.newTransferStructure.Fastresume))
			}
		})
	}
}

func TestTransferStructure_HandlePriority(t *testing.T) {

	type HandlePriorityCase struct {
		name                 string
		mustFail             bool
		newTransferStructure *TransferStructure
		expected             []int64
	}
	cases := []HandlePriorityCase{
		{
			name: "001 mustfail",
			newTransferStructure: &TransferStructure{
				Fastresume:  &qBittorrentStructures.QBittorrentFastresume{FilePriority: []int64{}},
				TorrentFile: &torrentStructures.Torrent{Info: &torrentStructures.TorrentInfo{}},
				ResumeItem: &utorrentStructs.ResumeItem{
					Prio: []byte{},
				},
			},
			mustFail: true,
			expected: []int64{0, 0, 1, 1, 1, 6, 6, 0},
		},
		{
			name: "002 check priotiry for v1 torrents",
			newTransferStructure: &TransferStructure{
				Fastresume:  &qBittorrentStructures.QBittorrentFastresume{FilePriority: []int64{}},
				TorrentFile: &torrentStructures.Torrent{Info: &torrentStructures.TorrentInfo{}},
				ResumeItem: &utorrentStructs.ResumeItem{
					Prio: []byte{
						byte(0),
						byte(128),
						byte(2),
						byte(5),
						byte(8),
						byte(9),
						byte(15),
						byte(127), // unexpected
					},
				},
			},
			expected: []int64{0, 0, 1, 1, 1, 6, 6, 0},
		},
		{
			name: "002 check priotiry for v2 torrents",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{FilePriority: []int64{}},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						FileTree: map[string]interface{}{},
					},
				},
				ResumeItem: &utorrentStructs.ResumeItem{
					Prio: []byte{
						byte(0),
						byte(128),
						byte(128),
						byte(128),
						byte(2),
						byte(128),
						byte(5),
						byte(128),
						byte(8),
						byte(128),
						byte(9),
						byte(128),
						byte(15),
						byte(128),
						byte(127), // unexpected
						byte(128),
					},
				},
			},
			expected: []int64{0, 0, 1, 1, 1, 6, 6, 0},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.newTransferStructure.HandlePriority()
			equal := reflect.DeepEqual(testCase.expected, testCase.newTransferStructure.Fastresume.FilePriority)
			if !equal && !testCase.mustFail {
				changes, err := diff.Diff(testCase.newTransferStructure.Fastresume.FilePriority, testCase.expected, diff.DiscardComplexOrigin())
				if err != nil {
					t.Error(err.Error())
				}
				t.Fatalf("Unexpected error: opts isn't equal:\n Got: %#v\n Expect %#v\n Diff: %v\n", testCase.newTransferStructure.Fastresume.FilePriority, testCase.expected, spew.Sdump(changes))
			} else if equal && testCase.mustFail {
				t.Fatalf("Unexpected error: structures are equal, but they shouldn't\n Got: %v\n", spew.Sdump(testCase.newTransferStructure.Fastresume.FilePriority))
			}
		})
	}
}

func TestTransferStructure_HandleTrackers(t *testing.T) {
	transferStructure := TransferStructure{
		Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
		ResumeItem: &utorrentStructs.ResumeItem{
			Trackers: []interface{}{
				"http://test1.org",
				"udp://test1.org",
				"http://test1.local",
				"udp://test1.local",
				[]interface{}{
					"http://test2.org:80",
					"udp://test2.org:8080",
					"http://test2.local:80",
					"udp://test2.local:8080",
					[]interface{}{
						"http://test3.org:80/somepath",
						"udp://test3.org:8080/somepath",
						"http://test3.local:80/somepath",
						"udp://test3.local:8080/somepath",
					},
				},
				[]interface{}{
					[]interface{}{
						"http://test4.org:80/",
						"udp://test4.org:8080/",
						"http://test4.local:80/",
						"udp://test4.local:8080/",
					},
				},
			},
		},
	}
	expect := [][]string{
		[]string{
			"http://test1.org", "udp://test1.org",
			"http://test2.org:80", "udp://test2.org:8080",
			"http://test3.org:80/somepath", "udp://test3.org:8080/somepath",
			"http://test4.org:80/", "udp://test4.org:8080/"},
		[]string{"http://test1.local", "udp://test1.local",
			"http://test2.local:80", "udp://test2.local:8080",
			"http://test3.local:80/somepath", "udp://test3.local:8080/somepath",
			"http://test4.local:80/", "udp://test4.local:8080/"},
	}
	transferStructure.HandleTrackers()
	if !reflect.DeepEqual(transferStructure.Fastresume.Trackers, expect) {
		t.Fatalf("Unexpected error: opts isn't equal:\n Got: %#v\n Expect %#v\n", transferStructure.Fastresume.Trackers, expect)
	}
}

func TestTransferStructure_HandleState(t *testing.T) {
	type HandleStateCase struct {
		name                 string
		mustFail             bool
		newTransferStructure *TransferStructure
		expected             *TransferStructure
	}
	cases := []HandleStateCase{
		{
			name: "001 Mustfail",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{},
							&torrentStructures.TorrentFile{},
							&torrentStructures.TorrentFile{},
						},
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
			},
			mustFail: true,
		},
		{
			name: "002 stopped resume",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Started: 0},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{},
							&torrentStructures.TorrentFile{},
							&torrentStructures.TorrentFile{},
						},
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Paused: 1, AutoManaged: 0},
			},
		},
		{
			name: "003 started resume",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{Started: 1},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Paused: 0, AutoManaged: 1},
			},
		},
		{
			name: "004 started resume with full downloaded files",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Started: 1,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
						},
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Paused: 0, AutoManaged: 1},
			},
		},
		{
			name: "005 started resume with parted downloaded files",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1, 1, 1, 6, 6},
				},
				ResumeItem: &utorrentStructs.ResumeItem{
					Started: 1,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
						},
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					Paused:       1,
					AutoManaged:  0,
					FilePriority: []int64{1, 0, 1, 1, 1, 6, 6},
				},
			},
		},
		{
			name: "006 started resume with parted downloaded files",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1, 1, 1, 6, 6},
				},
				ResumeItem: &utorrentStructs.ResumeItem{
					Started: 1,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
						},
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					Paused:       1,
					AutoManaged:  0,
					FilePriority: []int64{1, 0, 1, 1, 1, 6, 6},
				},
			},
		},
		{
			name: "007 started resume with full downloaded files",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1, 1, 1, 6, 6},
				},
				ResumeItem: &utorrentStructs.ResumeItem{
					Started: 1,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
						},
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					Paused:       1,
					AutoManaged:  0,
					FilePriority: []int64{1, 0, 1, 1, 1, 6, 6},
				},
			},
		},
		{
			name: "008 started resume with parted downloaded files",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					FilePriority: []int64{1, 0, 1, 1, 1, 6, 6},
				},
				ResumeItem: &utorrentStructs.ResumeItem{
					Started: 1,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{
						Files: []*torrentStructures.TorrentFile{
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
							&torrentStructures.TorrentFile{Path: []string{`/`}},
						},
					},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{
					Paused:       1,
					AutoManaged:  0,
					FilePriority: []int64{1, 0, 1, 1, 1, 6, 6},
				},
			},
		},
		{
			name: "009 started resume without files",
			newTransferStructure: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{},
				ResumeItem: &utorrentStructs.ResumeItem{
					Started: 0,
				},
				TorrentFile: &torrentStructures.Torrent{
					Info: &torrentStructures.TorrentInfo{},
				},
			},
			expected: &TransferStructure{
				Fastresume: &qBittorrentStructures.QBittorrentFastresume{Paused: 1, AutoManaged: 0},
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.newTransferStructure.HandleState()
			equal := reflect.DeepEqual(testCase.expected.Fastresume, testCase.newTransferStructure.Fastresume)
			if !equal && !testCase.mustFail {
				changes, err := diff.Diff(testCase.newTransferStructure.Fastresume, testCase.expected.Fastresume, diff.DiscardComplexOrigin())
				if err != nil {
					t.Error(err.Error())
				}
				t.Fatalf("Unexpected error: opts isn't equal:\n Got: %#v\n Expect %#v\n Diff: %v\n", testCase.newTransferStructure.Fastresume.Pieces, testCase.expected.Fastresume.Pieces, spew.Sdump(changes))
			} else if equal && testCase.mustFail {
				t.Fatalf("Unexpected error: structures are equal, but they shouldn't\n Got: %v\n", spew.Sdump(testCase.newTransferStructure.Fastresume))
			}
		})
	}

}
