package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// fungsi ini akan membuat program secara otomatis
// mengupdate domain list setiap 1 jam.
func domainsFetch() {

	for {

		// pasang user agent
		req, err := http.NewRequest("GET", "https://raw.githubusercontent.com/lepasid/blocklist/main/domains", nil)
		if err != nil {
			fmt.Println(" [req] ", err)
			return
		}
		req.Header.Set("User-Agent", usrAgent)

		// lakukan fetch untuk ambil domain list
		getData, err := h1Client.Do(req)
		if err != nil {
			fmt.Println(" [getData] ", err)
			return
		}

		// buat file baru di dalam memory
		createFile, err := memFS.Create(kominfodDir)
		if err != nil {
			fmt.Println(" [createFile] ", err)

			return
		}

		// tulis data hasil fetch ke file yang sudah dibuat
		writeFile, err := io.Copy(createFile, getData.Body)
		if err != nil {
			fmt.Println(" [writeHTML] ", err)
			return
		}

		// lakukan print sebagi indikasi bahwa data sudah berhasil diproses
		sizeinfo := fmt.Sprintf("Downloading %v KB | %v MB", (writeFile / Kilobyte), (writeFile / Megabyte))
		fmt.Println(sizeinfo)

		// tutup io setelah selesai menulis data
		if err := createFile.Close(); err != nil {
			fmt.Println(" [createFile.Close()] ", err)
			return
		}

		// tutup http body setelah digunakan
		getData.Body.Close()

		// lakukan sleep selama 1 jam
		time.Sleep(1 * time.Hour)
	}
}
