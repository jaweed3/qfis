package db

import (
	"fmt"
	"math"
	"math/rand"

	"qfis/internal/models"
)

var seedMerchants = []struct {
	name       string
	mcc        string
	npwp       bool
	flagged     bool
	suspicious bool
}{
	{name: "Pulsa Murah 24Jam", mcc: "5999", npwp: false, flagged: true, suspicious: true},
	{name: "Token Digital Express", mcc: "5999", npwp: false, flagged: true, suspicious: true},
	{name: "Reload Gaming Center", mcc: "7993", npwp: true, flagged: true, suspicious: true},
	{name: "Top Up Instan", mcc: "5999", npwp: false, flagged: true, suspicious: true},
	{name: "Isi Saldo Gaming", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Voucher Gaming Indo", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Pulsa Cepat Sakti", mcc: "5999", npwp: false, flagged: true, suspicious: true},
	{name: "Refill Akun Pro", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Chip Poker Online", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Depo Gaming Cepat", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Warung Kopi Barokah", mcc: "5812", npwp: true, flagged: false, suspicious: false},
	{name: "Alfamart Gontor", mcc: "5411", npwp: true, flagged: false, suspicious: false},
	{name: "Toko Sembako Ibu", mcc: "5411", npwp: true, flagged: false, suspicious: false},
	{name: "Bakso Mas Joko", mcc: "5812", npwp: true, flagged: false, suspicious: false},
	{name: "Laundry Express", mcc: "7216", npwp: true, flagged: false, suspicious: false},
	{name: "Bengkel Motor Jaya", mcc: "7538", npwp: true, flagged: false, suspicious: false},
	{name: "Salon Cantik", mcc: "7297", npwp: true, flagged: false, suspicious: false},
	{name: "Martabak Bang Ali", mcc: "5812", npwp: true, flagged: false, suspicious: false},
	{name: "Minimarket Sejahtera", mcc: "5411", npwp: true, flagged: false, suspicious: false},
	{name: "Photo Copy Cepat", mcc: "7338", npwp: true, flagged: false, suspicious: false},
	{name: "Token Pulsa Digital", mcc: "5999", npwp: false, flagged: true, suspicious: true},
	{name: "Isi Ulang Game", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Voucher Diskon", mcc: "5999", npwp: false, flagged: true, suspicious: true},
	{name: "Top Up Game Murah", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Pulsa Elektrik", mcc: "5999", npwp: false, flagged: true, suspicious: true},
	{name: "Nasi Goreng Mantap", mcc: "5812", npwp: true, flagged: false, suspicious: false},
	{name: "Ayam Geprek Juara", mcc: "5812", npwp: true, flagged: false, suspicious: false},
	{name: "Es Teh Indonesia", mcc: "5812", npwp: true, flagged: false, suspicious: false},
	{name: "Toko Roti Fresh", mcc: "5462", npwp: true, flagged: false, suspicious: false},
	{name: "Kerupuk Mawar", mcc: "5411", npwp: true, flagged: false, suspicious: false},
	{name: "Agen Pulsa Cepat", mcc: "5999", npwp: false, flagged: true, suspicious: true},
	{name: "Game Station ID", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Mi Ayam Pakde", mcc: "5812", npwp: true, flagged: false, suspicious: false},
	{name: "Soto Ayam Bu Tini", mcc: "5812", npwp: true, flagged: false, suspicious: false},
	{name: "Pecel Lele Lela", mcc: "5812", npwp: true, flagged: false, suspicious: false},
	{name: "Top Up Mobile Legend", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Voucher Game Online", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Konter Pulsa Sam", mcc: "5999", npwp: true, flagged: false, suspicious: false},
	{name: "Toko Elektronik Makmur", mcc: "5999", npwp: true, flagged: false, suspicious: false},
	{name: "Buku Stationery", mcc: "5942", npwp: true, flagged: false, suspicious: false},
	{name: "Depo Chip Poker", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Isi Chip Online", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Saldo Game Cepat", mcc: "7993", npwp: false, flagged: true, suspicious: true},
	{name: "Cemilan Enak", mcc: "5411", npwp: true, flagged: false, suspicious: false},
	{name: "Tahu Bulat Pak Jay", mcc: "5812", npwp: true, flagged: false, suspicious: false},
}

func nameToID(name string) string {
	h := 0
	for _, c := range name {
		h = h*31 + int(c)
	}
	rng := rand.New(rand.NewSource(int64(h)))
	num := rng.Intn(9000) + 1000
	letter := string(rune('A' + rng.Intn(26)))
	return fmt.Sprintf("QRIS-%d%s", num, letter)
}

func Seed() error {
	var merchantIDs []string

	for i, s := range seedMerchants {
		merchantID := nameToID(s.name)
		merchantIDs = append(merchantIDs, merchantID)

		rng := rand.New(rand.NewSource(int64(i * 42)))
		x := 0.05 + rng.Float64()*0.9
		y := 0.05 + rng.Float64()*0.9

		m := &models.Merchant{
			ID:      merchantID,
			Name:    s.name,
			MCC:     s.mcc,
			NPWP:    s.npwp,
			Flagged: s.flagged,
			X:       math.Round(x*100) / 100,
			Y:       math.Round(y*100) / 100,
		}
		if err := UpsertMerchant(m); err != nil {
			return fmt.Errorf("seed merchant %s: %w", s.name, err)
		}
	}

	rng := rand.New(rand.NewSource(42))
	for i, s := range seedMerchants {
		if !s.suspicious {
			continue
		}
		merchantID := merchantIDs[i]
		score := math.Round((65+rng.Float64()*30)*10) / 10

		mockReport := &models.Report{
			ID:           fmt.Sprintf("RPT-SEED-%03d", i+1),
			MerchantID:   merchantID,
			Note:         fmt.Sprintf("Laporan crowdsource: merchant %s terindikasi aktivitas mencurigakan", s.name),
			Score:        score,
			Status:       "queued",
		}
		if err := InsertReport(mockReport); err != nil {
			return fmt.Errorf("seed report: %w", err)
		}
	}

	return nil
}
