package tools

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"time"
)

func InitAutoMakeTLSConfigWithArgs(cert, key string) *tls.Config {
	certVal, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatalf("Cert load failed: %v", err)
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{certVal},
	}

	return conf
}

func InitAutoMakeTLSConfig() *tls.Config {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{"IMAIL"},
		OrganizationalUnit: []string{"IMAIL ORG Unit"},
		CommonName:         "IMAIL",
	}
	rootTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	makeCert, _ := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &pk.PublicKey, pk)

	privBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pk),
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: makeCert,
	}

	cert, err := tls.X509KeyPair(pem.EncodeToMemory(certBlock), pem.EncodeToMemory(privBlock))

	//------- demo start ------

	// 	var localhostCert = []byte(`-----BEGIN CERTIFICATE-----
	// MIIFkzCCA3ugAwIBAgIUQvhoyGmvPHq8q6BHrygu4dPp0CkwDQYJKoZIhvcNAQEL
	// BQAwWTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
	// GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDESMBAGA1UEAwwJbG9jYWxob3N0MB4X
	// DTIwMDUyMTE2MzI1NVoXDTMwMDUxOTE2MzI1NVowWTELMAkGA1UEBhMCQVUxEzAR
	// BgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5
	// IEx0ZDESMBAGA1UEAwwJbG9jYWxob3N0MIICIjANBgkqhkiG9w0BAQEFAAOCAg8A
	// MIICCgKCAgEAk773plyfK4u2uIIZ6H7vEnTb5qJT6R/KCY9yniRvCFV+jCrISAs9
	// 0pgU+/P8iePnZRGbRCGGt1B+1/JAVLIYFZuawILHNs4yWKAwh0uNpR1Pec8v7vpq
	// NpdUzXKQKIqFynSkcLA8c2DOZwuhwVc8rZw50yY3r4i4Vxf0AARGXapnBfy6WerR
	// /6xT7y/OcK8+8aOirDQ9P6WlvZ0ynZKi5q2o1eEVypT2us9r+HsCYosKEEAnjzjJ
	// wP5rvredxUqb7OupIkgA4Nq80+4tqGGQfWetmoi3zXRhKpijKjgxBOYEqSUWm9ws
	// /aC91Iy5RawyTB0W064z75OgfuI5GwFUbyLD0YVN4DLSAI79GUfvc8NeLEXpQvYq
	// +f8P+O1Hbv2AQ28IdbyQrNefB+/WgjeTvXLploNlUihVhpmLpptqnauw/DY5Ix51
	// w60lHIZ6esNOmMQB+/z/IY5gpmuo66yH8aSCPSYBFxQebB7NMqYGOS9nXx62/Bn1
	// OUVXtdtrhfbbdQW6zMZjka0t8m83fnGw3ISyBK2NNnSzOgycu0ChsW6sk7lKyeWa
	// 85eJGsQWIhkOeF9v9GAIH/qsrgVpToVC9Krbk+/gqYIYF330tHQrzp6M6LiG5OY1
	// P7grUBovN2ZFt10B97HxWKa2f/8t9sfHZuKbfLSFbDsyI2JyNDh+Vk0CAwEAAaNT
	// MFEwHQYDVR0OBBYEFOLdIQUr3gDQF5YBor75mlnCdKngMB8GA1UdIwQYMBaAFOLd
	// IQUr3gDQF5YBor75mlnCdKngMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQEL
	// BQADggIBAGddhQMVMZ14TY7bU8CMuc9IrXUwxp59QfqpcXCA2pHc2VOWkylv2dH7
	// ta6KooPMKwJ61d+coYPK1zMUvNHHJCYVpVK0r+IGzs8mzg91JJpX2gV5moJqNXvd
	// Fy6heQJuAvzbb0Tfsv8KN7U8zg/ovpS7MbY+8mRJTQINn2pCzt2y2C7EftLK36x0
	// KeBWqyXofBJoMy03VfCRqQlWK7VPqxluAbkH+bzji1g/BTkoCKzOitAbjS5lT3sk
	// oCrF9N6AcjpFOH2ZZmTO4cZ6TSWfrb/9OWFXl0TNR9+x5c/bUEKoGeSMV1YT1SlK
	// TNFMUlq0sPRgaITotRdcptc045M6KF777QVbrYm/VH1T3pwPGYu2kUdYHcteyX9P
	// 8aRG4xsPGQ6DD7YjBFsif2fxlR3nQ+J/l/+eXHO4C+eRbxi15Z2NjwVjYpxZlUOq
	// HD96v516JkMJ63awbY+HkYdEUBKqR55tzcvNWnnfiboVmIecjAjoV4zStwDIti9u
	// 14IgdqqAbnx0ALbUWnvfFloLdCzPPQhgLHpTeRSEDPljJWX8rmy8iQtRb0FWYQ3z
	// A2wsUyutzK19nt4hjVrTX0At9ku3gMmViXFlbvyA1Y4TuhdUYqJauMBrWKl2ybDW
	// yhdKg/V3yTwgBUtb3QO4m1khNQjQLuPFVxULGEA38Y5dXSONsYnt
	// -----END CERTIFICATE-----`)

	// 	var localhostKey = []byte(`-----BEGIN PRIVATE KEY-----
	// MIIJQgIBADANBgkqhkiG9w0BAQEFAASCCSwwggkoAgEAAoICAQCTvvemXJ8ri7a4
	// ghnofu8SdNvmolPpH8oJj3KeJG8IVX6MKshICz3SmBT78/yJ4+dlEZtEIYa3UH7X
	// 8kBUshgVm5rAgsc2zjJYoDCHS42lHU95zy/u+mo2l1TNcpAoioXKdKRwsDxzYM5n
	// C6HBVzytnDnTJjeviLhXF/QABEZdqmcF/LpZ6tH/rFPvL85wrz7xo6KsND0/paW9
	// nTKdkqLmrajV4RXKlPa6z2v4ewJiiwoQQCePOMnA/mu+t53FSpvs66kiSADg2rzT
	// 7i2oYZB9Z62aiLfNdGEqmKMqODEE5gSpJRab3Cz9oL3UjLlFrDJMHRbTrjPvk6B+
	// 4jkbAVRvIsPRhU3gMtIAjv0ZR+9zw14sRelC9ir5/w/47Udu/YBDbwh1vJCs158H
	// 79aCN5O9cumWg2VSKFWGmYumm2qdq7D8NjkjHnXDrSUchnp6w06YxAH7/P8hjmCm
	// a6jrrIfxpII9JgEXFB5sHs0ypgY5L2dfHrb8GfU5RVe122uF9tt1BbrMxmORrS3y
	// bzd+cbDchLIErY02dLM6DJy7QKGxbqyTuUrJ5Zrzl4kaxBYiGQ54X2/0YAgf+qyu
	// BWlOhUL0qtuT7+CpghgXffS0dCvOnozouIbk5jU/uCtQGi83ZkW3XQH3sfFYprZ/
	// /y32x8dm4pt8tIVsOzIjYnI0OH5WTQIDAQABAoICADBPw788jje5CdivgjVKPHa2
	// i6mQ7wtN/8y8gWhA1aXN/wFqg+867c5NOJ9imvOj+GhOJ41RwTF0OuX2Kx8G1WVL
	// aoEEwoujRUdBqlyzUe/p87ELFMt6Svzq4yoDCiyXj0QyfAr1Ne8sepGrdgs4sXi7
	// mxT2bEMT2+Nuy7StsSyzqdiFWZJJfL2z5gZShZjHVTfCoFDbDCQh0F5+Zqyr5GS1
	// 6H13ip6hs0RGyzGHV7JNcM77i3QDx8U57JWCiS6YRQBl1vqEvPTJ0fEi8v8aWBsJ
	// qfTcO+4M3jEFlGUb1ruZU3DT1d7FUljlFO3JzlOACTpmUK6LSiRPC64x3yZ7etYV
	// QGStTdjdJ5+nE3CPR/ig27JLrwvrpR6LUKs4Dg13g/cQmhpq30a4UxV+y8cOgR6g
	// 13YFOtZto2xR+53aP6KMbWhmgMp21gqxS+b/5HoEfKCdRR1oLYTVdIxt4zuKlfQP
	// pTjyFDPA257VqYy+e+wB/0cFcPG4RaKONf9HShlWAulriS/QcoOlE/5xF74QnmTn
	// YAYNyfble/V2EZyd2doU7jJbhwWfWaXiCMOO8mJc+pGs4DsGsXvQmXlawyElNWes
	// wJfxsy4QOcMV54+R/wxB+5hxffUDxlRWUsqVN+p3/xc9fEuK+GzuH+BuI01YQsw/
	// laBzOTJthDbn6BCxdCeBAoIBAQDEO1hDM4ZZMYnErXWf/jik9EZFzOJFdz7g+eHm
	// YifFiKM09LYu4UNVY+Y1btHBLwhrDotpmHl/Zi3LYZQscWkrUbhXzPN6JIw98mZ/
	// tFzllI3Ioqf0HLrm1QpG2l7Xf8HT+d3atEOtgLQFYehjsFmmJtE1VsRWM1kySLlG
	// 11bQkXAlv7ZQ13BodQ5kNM3KLvkGPxCNtC9VQx3Em+t/eIZOe0Nb2fpYzY/lH1mF
	// rFhj6xf+LFdMseebOCQT27bzzlDrvWobQSQHqflFkMj86q/8I8RUAPcRz5s43YdO
	// Q+Dx2uJQtNBAEQVoS9v1HgBg6LieDt0ZytDETR5G3028dyaxAoIBAQDAvxEwfQu2
	// TxpeYQltHU/xRz3blpazgkXT6W4OT43rYI0tqdLxIFRSTnZap9cjzCszH10KjAg5
	// AQDd7wN6l0mGg0iyL0xjWX0cT38+wiz0RdgeHTxRk208qTyw6Xuh3KX2yryHLtf5
	// s3z5zkTJmj7XXOC2OVsiQcIFPhVXO3d38rm0xvzT5FZQH3a5rkpks1mqTZ4dyvim
	// p6vey4ZXdUnROiNzqtqbgSLbyS7vKj5/fXbkgKh8GJLNV4LMD6jo2FRN/LsEZKes
	// pxWNMsHBkv5eRfHNBVZuUMKFenN6ojV2GFG7bvLYD8Z9sja8AuBCaMr1CgHD8kd5
	// +A5+53Iva8hdAoIBAFU+BlBi8IiMaXFjfIY80/RsHJ6zqtNMQqdORWBj4S0A9wzJ
	// BN8Ggc51MAqkEkAeI0UGM29yicza4SfJQqmvtmTYAgE6CcZUXAuI4he1jOk6CAFR
	// Dy6O0G33u5gdwjdQyy0/DK21wvR6xTjVWDL952Oy1wyZnX5oneWnC70HTDIcC6CK
	// UDN78tudhdvnyEF8+DZLbPBxhmI+Xo8KwFlGTOmIyDD9Vq/+0/RPEv9rZ5Y4CNsj
	// /eRWH+sgjyOFPUtZo3NUe+RM/s7JenxKsdSUSlB4ZQ+sv6cgDSi9qspH2E6Xq9ot
	// QY2jFztAQNOQ7c8rKQ+YG1nZ7ahoa6+Tz1wAUnECggEAFVTP/TLJmgqVG37XwTiu
	// QUCmKug2k3VGbxZ1dKX/Sd5soXIbA06VpmpClPPgTnjpCwZckK9AtbZTtzwdgXK+
	// 02EyKW4soQ4lV33A0lxBB2O3cFXB+DE9tKnyKo4cfaRixbZYOQnJIzxnB2p5mGo2
	// rDT+NYyRdnAanePqDrZpGWBGhyhCkNzDZKimxhPw7cYflUZzyk5NSHxj/AtAOeuk
	// GMC7bbCp8u3Ows44IIXnVsq23sESZHF/xbP6qMTO574RTnQ66liNagEv1Gmaoea3
	// ug05nnwJvbm4XXdY0mijTAeS/BBiVeEhEYYoopQa556bX5UU7u+gU3JNgGPy8iaW
	// jQKCAQEAp16lci8FkF9rZXSf5/yOqAMhbBec1F/5X/NQ/gZNw9dDG0AEkBOJQpfX
	// dczmNzaMSt5wmZ+qIlu4nxRiMOaWh5LLntncQoxuAs+sCtZ9bK2c19Urg5WJ615R
	// d6OWtKINyuVosvlGzquht+ZnejJAgr1XsgF9cCxZonecwYQRlBvOjMRidCTpjzCu
	// 6SEEg/JyiauHq6wZjbz20fXkdD+P8PIV1ZnyUIakDgI7kY0AQHdKh4PSMvDoFpIw
	// TXU5YrNA8ao1B6CFdyjmLzoY2C9d9SDQTXMX8f8f3GUo9gZ0IzSIFVGFpsKBU0QM
	// hBgHM6A0WJC9MO3aAKRBcp48y6DXNA==
	// -----END PRIVATE KEY-----`)

	// 	cert, err := tls.X509KeyPair(localhostCert, localhostKey)

	//------- demo end ------

	if err != nil {
		log.Fatalf("Cert load failed: %v", err)
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return conf
}
