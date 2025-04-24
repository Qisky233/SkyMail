package model

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/yanyiwu/gojieba"
)

var (
	globalJieba *gojieba.Jieba
	jiebaOnce   sync.Once
)

func getGlobalJieba() *gojieba.Jieba {
	jiebaOnce.Do(func() {
		globalJieba = gojieba.NewJieba()
	})
	return globalJieba
}

type SpamFilter struct {
	WordDict  map[string][2]int // [ham_count, spam_count]
	HamTotal  int
	SpamTotal int
	Jieba     *gojieba.Jieba
}

func NewSpamFilter(modelPath string) (*SpamFilter, error) {
	sf := &SpamFilter{
		WordDict: make(map[string][2]int),
		Jieba:    getGlobalJieba(),
	}

	if err := sf.loadModel(modelPath); err != nil {
		return nil, fmt.Errorf("failed to load model: %v", err)
	}

	return sf, nil
}

func (sf *SpamFilter) loadModel(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check magic number
	magic := make([]byte, 4)
	if _, err := file.Read(magic); err != nil {
		return err
	}
	if string(magic) != "GONB" {
		return fmt.Errorf("invalid model file format")
	}

	// Read version
	var version uint32
	if err := binary.Read(file, binary.LittleEndian, &version); err != nil {
		return err
	}
	if version != 1 {
		return fmt.Errorf("unsupported model version: %d", version)
	}

	// Read totals
	var hamTotal, spamTotal uint32
	if err := binary.Read(file, binary.LittleEndian, &hamTotal); err != nil {
		return err
	}
	if err := binary.Read(file, binary.LittleEndian, &spamTotal); err != nil {
		return err
	}
	sf.HamTotal = int(hamTotal)
	sf.SpamTotal = int(spamTotal)

	// Read dictionary size
	var dictSize uint32
	if err := binary.Read(file, binary.LittleEndian, &dictSize); err != nil {
		return err
	}

	for i := 0; i < int(dictSize); i++ {
		// Read word length
		var wordLen uint32
		if err := binary.Read(file, binary.LittleEndian, &wordLen); err != nil {
			return err
		}

		// Read word
		wordBytes := make([]byte, wordLen)
		if _, err := file.Read(wordBytes); err != nil {
			return err
		}
		word := string(wordBytes)

		// Read counts
		var hamCount, spamCount uint32
		if err := binary.Read(file, binary.LittleEndian, &hamCount); err != nil {
			return err
		}
		if err := binary.Read(file, binary.LittleEndian, &spamCount); err != nil {
			return err
		}
		sf.WordDict[word] = [2]int{int(hamCount), int(spamCount)}
	}

	return nil
}

func (sf *SpamFilter) Preprocess(text string) []string {
	var chineseText strings.Builder
	for _, r := range text {
		if r >= 0x4E00 && r <= 0x9FD5 {
			chineseText.WriteRune(r)
		}
	}

	words := sf.Jieba.Cut(chineseText.String(), true)
	var result []string
	for _, word := range words {
		if word != "" {
			result = append(result, word)
		}
	}
	return result
}

func (sf *SpamFilter) Predict(text string, threshold float64) bool {
	words := sf.Preprocess(text)

	hamProb := math.Log(float64(sf.HamTotal) / float64(sf.HamTotal+sf.SpamTotal))
	spamProb := math.Log(float64(sf.SpamTotal) / float64(sf.HamTotal+sf.SpamTotal))

	for _, word := range words {
		counts, exists := sf.WordDict[word]
		if !exists {
			hamProb += math.Log(1.0 / float64(sf.HamTotal+1))
			spamProb += math.Log(1.0 / float64(sf.SpamTotal+1))
			continue
		}

		if counts[0] > 0 {
			hamProb += math.Log(float64(counts[0]) / float64(sf.HamTotal))
		} else {
			hamProb += math.Log(1.0 / float64(sf.HamTotal+1))
		}

		if counts[1] > 0 {
			spamProb += math.Log(float64(counts[1]) / float64(sf.SpamTotal))
		} else {
			spamProb += math.Log(1.0 / float64(sf.SpamTotal+1))
		}
	}

	return spamProb+threshold > hamProb
}

func (sf *SpamFilter) Close() {
	// 不需要释放 globalJieba
}
