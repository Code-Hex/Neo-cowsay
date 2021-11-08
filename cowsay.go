package cowsay

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Say to return cowsay string.
func Say(phrase string, options ...Option) (string, error) {
	cow, err := New(options...)
	if err != nil {
		return "", err
	}
	return cow.Say(phrase)
}

// LocationType indicates the type of COWPATH.
type LocationType int

const (
	// InBinary indicates the COWPATH in binary.
	InBinary LocationType = iota

	// InDirectory indicates the COWPATH in your directory.
	InDirectory
)

// CowPath is information of the COWPATH.
type CowPath struct {
	// Name is name of the COWPATH.
	// If you specified `COWPATH=/foo/bar`, Name is `/foo/bar`.
	Name string
	// CowFiles are name of the cowfile which are trimmed ".cow" suffix.
	CowFiles []string
	// LocationType is the type of COWPATH
	LocationType LocationType
}

// Lookup will look for the target cowfile in the specified path.
// If it exists, it returns the cowfile information and true value.
func (c *CowPath) Lookup(target string) (*CowFile, bool) {
	for _, cowfile := range c.CowFiles {
		if cowfile == target {
			return &CowFile{
				Name:         cowfile,
				BasePath:     c.Name,
				LocationType: c.LocationType,
			}, true
		}
	}
	return nil, false
}

// CowFile is information of the cowfile.
type CowFile struct {
	// Name is name of the cowfile.
	Name string
	// BasePath is the path which the cowpath is in.
	BasePath string
	// LocationType is the type of COWPATH
	LocationType LocationType
}

// ReadAll reads the cowfile content.
// If LocationType is InBinary, the file read from binary.
// otherwise reads from file system.
func (c *CowFile) ReadAll() ([]byte, error) {
	joinedPath := filepath.Join(c.BasePath, c.Name+".cow")
	if c.LocationType == InBinary {
		return Asset(joinedPath)
	}
	return ioutil.ReadFile(joinedPath)
}

// Cows to get list of cows
func Cows() ([]*CowPath, error) {
	cowPaths, err := cowsFromCowPath()
	if err != nil {
		return nil, err
	}
	cowPaths = append(cowPaths, &CowPath{
		Name:         "cows",
		CowFiles:     CowsInBinary(),
		LocationType: InBinary,
	})
	return cowPaths, nil
}

func cowsFromCowPath() ([]*CowPath, error) {
	cowPaths := make([]*CowPath, 0)
	cowPath := os.Getenv("COWPATH")
	if cowPath == "" {
		return cowPaths, nil
	}
	paths := splitPath(cowPath)
	for _, path := range paths {
		dirEntries, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		path := &CowPath{
			Name:         path,
			CowFiles:     []string{},
			LocationType: InDirectory,
		}
		for _, entry := range dirEntries {
			name := entry.Name()
			if strings.HasSuffix(name, ".cow") {
				name = strings.TrimSuffix(name, ".cow")
				path.CowFiles = append(path.CowFiles, name)
			}
		}
		sort.Strings(path.CowFiles)
		cowPaths = append(cowPaths, path)
	}
	return cowPaths, nil
}

// GetCow to get cow's ascii art
func (cow *Cow) GetCow() (string, error) {
	src, err := cow.typ.ReadAll()
	if err != nil {
		return "", err
	}

	if len(cow.eyes) > 2 {
		cow.eyes = cow.eyes[0:2]
	}

	if len(cow.tongue) > 2 {
		cow.tongue = cow.tongue[0:2]
	}

	r := strings.NewReplacer(
		"\\\\", "\\",
		"\\@", "@",
		"\\$", "$",
		"$eyes", cow.eyes,
		"${eyes}", cow.eyes,
		"$tongue", cow.tongue,
		"${tongue}", cow.tongue,
		"$thoughts", string(cow.thoughts),
		"${thoughts}", string(cow.thoughts),
	)
	newsrc := r.Replace(string(src))
	separate := strings.Split(newsrc, "\n")
	mow := make([]string, 0, len(separate))
	for _, line := range separate {
		if strings.Contains(line, "$the_cow = <<EOC") || strings.HasPrefix(line, "##") {
			continue
		}

		if strings.HasPrefix(line, "EOC") {
			break
		}

		mow = append(mow, line)
	}
	return strings.Join(mow, "\n"), nil
}
