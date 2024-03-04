# go-kominfod

## Daftar Isi

1. [Pengenalan](https://github.com/galpt/go-kominfod?tab=readme-ov-file#pengenalan)
2. [Apa yang dapat diharapkan?](https://github.com/galpt/go-kominfod?tab=readme-ov-file#what-to-expect)
3. [Cara kerja](https://github.com/galpt/go-kominfod?tab=readme-ov-file#how-it-works)
4. [Cara compile kode](https://github.com/galpt/go-kominfod?tab=readme-ov-file#how-to-compile-the-code)
5. [Lihat go-kominfod secara langsung](https://github.com/galpt/go-kominfod?tab=readme-ov-file#see-go-kominfod-in-action)
6. [Credits](https://github.com/galpt/go-kominfod?tab=readme-ov-file#credits)

* * *

## Pengenalan
#### [:arrow_up: Kembali ke Daftar Isi](https://github.com/galpt/go-kominfod?tab=readme-ov-file#daftar-isi)

`go-kominfod` adalah sebuah implementasi *low latency* yang diadaptasi dari proyek [kominfod](https://github.com/lepasid/kominfod), ditulis dalam bahasa Go.

Kata *"low latency"* yang dimaksud disini adalah kemampuan program ini untuk memproses permintaan dari pengguna dan memberikan respon balik dalam waktu sesingkat mungkin.

`go-kominfod` menggunakan beberapa library yang sudah diuji oleh banyak developer lainnya seperti [Gin](https://github.com/gin-gonic/gin) dan [Afero](https://github.com/spf13/afero) yang membuat `go-kominfod` menjadi lebih kokoh untuk dijalankan tanpa henti, baik itu di testing maupun production environment.

* * *

## Apa yang dapat diharapkan?
#### [:arrow_up: Kembali ke Daftar Isi](https://github.com/galpt/go-kominfod?tab=readme-ov-file#daftar-isi)

Ada beberapa hal yang dapat diharapkan dari implementasi ini, yaitu:
1. Dirancang untuk low latency. Rerata request diproses dalam waktu 100ms - 200ms.
2. Tidak melakukan proses baca-tulis ke disk sama sekali.
3. Secara teori dapat memproses jutaan request dengan cepat karena menggunakan library yang sudah diuji seperti [Gin](https://github.com/gin-gonic/gin) dan [Afero](https://github.com/spf13/afero).

* * *

## Cara kerja
#### [:arrow_up: Kembali ke Daftar Isi](https://github.com/galpt/go-kominfod?tab=readme-ov-file#daftar-isi)

Cara kerja `go-kominfod` jika dijelaskan setiap langkah, yaitu:
1. Saat program baru dijalankan, program akan mengunduh file [domains](https://raw.githubusercontent.com/lepasid/blocklist/main/domains) lalu disimpan di memory.
2. Ketika pengguna mengirim request ke `/kominfod`, secara default program akan mengecek apakah `google.com` diblokir atau tidak.
3. Jika pengguna menggunakan format `/kominfod?domain=`, maka program akan mengecek apakah domain tersebut diblokir atau tidak.
4. Jika pengguna menggunakan format `/kominfod?domain=google.com,facebook.com`, maka program akan mengecek seluruh domain tersebut apakah diblokir atau tidak.

* * *

## Cara compile kode
#### [:arrow_up: Kembali ke Daftar Isi](https://github.com/galpt/go-kominfod?tab=readme-ov-file#daftar-isi)

1. Download dan install [bahasa pemrograman Go](https://go.dev/).
2. Download repository ini dan extract ke sebuah folder baru yang masih kosong.
3. Compile kode dengan cara berikut:

```yaml
$ cd ./go-kominfod
$ go mod tidy
$ go build
```

> [!IMPORTANT]
> 1. Jalankan file yang sudah di-compile dengan `sudo` (untuk Linux) atau `Run as administrator` (untuk Windows) jika diperlukan.
> 2. Izinkanlah jika ada pop-up firewall untuk memberi izin program mengakses port `7777` (default).
> 3. Ubah `CertFilePath` dan `KeyFilePath` pada file `vars.go` sesuai lokasi sertifikat SSL yang ingin digunakan.

* * *

## Lihat `go-kominfod` secara langsung
#### [:arrow_up: Kembali ke Daftar Isi](https://github.com/galpt/go-kominfod?tab=readme-ov-file#daftar-isi)

Anda dapat melihat `go-kominfod` yang dijalankan di server kami dengan mengakses link berikut:

https://net.0ms.dev:7777/kominfod?domain=google.com,facebook.com

* * *

## Credits
#### [:arrow_up: Kembali ke Daftar Isi](https://github.com/galpt/go-kominfod?tab=readme-ov-file#daftar-isi)

Implementasi ini memungkinkan untuk dibuat dengan menggunakan hal-hal lain yang disediakan oleh para developer dan/atau perusahaan yang disebut disini.

Semua kredit dan hak cipta diberikan kepada masing-masing pemilik.
