package logman_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"moul.io/logman"
	"moul.io/u"
)

/* FIXME: various tests
 *
 * - input validation
 * - fs manipulation (empty dir, existing dir)
 * - GC safety (keep non-log files)
 * - concurrency
 */

func TestLogfile(t *testing.T) {
	// setup volatile directory for the test
	tempdir, err := ioutil.TempDir("", "logutil-file")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)

	// init logman
	manager := logman.Manager{
		Path:     tempdir,
		MaxFiles: 10,
	}

	invalidManager := logman.Manager{
		Path:     filepath.Join(tempdir, "doesnotexist"),
		MaxFiles: 10,
	}

	// check loading log files from an invalid directory
	{
		files, err := invalidManager.Files()
		require.Error(t, err)
		require.Nil(t, files)
	}

	// check loading files from empty valid directory
	{
		files, err := manager.Files()
		require.NoError(t, err)
		require.Empty(t, files)
	}

	// create dummy files
	{
		dummyNames := []string{
			"2021-05-25T21-12-02.650.log",
			"cli.info-2021-05-25T21-12-02.aaa.log",
			"blah.log",
		}
		for _, name := range dummyNames {
			f, err := os.Create(filepath.Join(tempdir, name))
			require.NoError(t, err)
			err = f.Close()
			require.NoError(t, err)
		}
	}

	// check loading files from valid directory with only dummy files
	{
		files, err := manager.Files()
		require.NoError(t, err)
		require.Empty(t, files)
	}

	// create a first logger of kind-1
	{
		writer, err := manager.New("kind-1")
		require.NoError(t, err)
		require.NotNil(t, writer)
		_, err = writer.Write([]byte("blah\n"))
		require.NoError(t, err)
		err = writer.Close()
		require.NoError(t, err)
	}

	// check loading files from the directory, should have one now
	{
		files, err := manager.Files()
		require.NoError(t, err)
		require.Len(t, files, 1)
		require.Equal(t, filepath.Dir(files[0].Path), tempdir)
		require.NotEmpty(t, files[0].Name)
		require.NotEmpty(t, files[0].Path)
		require.True(t, u.FileExists(files[0].Path))
		require.True(t, files[0].Latest)
		require.Equal(t, files[0].Name, "kind-1")
	}

	// create a second logger of kind-1
	{
		time.Sleep(time.Second)
		writer, err := manager.New("kind-1")
		require.NoError(t, err)
		require.NotNil(t, writer)
		_, err = writer.Write([]byte("blah blah\n"))
		require.NoError(t, err)
		err = writer.Close()
		require.NoError(t, err)
	}

	// check loading files from the directory, should have two now
	{
		files, err := manager.Files()
		require.NoError(t, err)
		require.Len(t, files, 2)
		for _, file := range files {
			require.Equal(t, filepath.Dir(file.Path), tempdir)
			require.NotEmpty(t, file.Name)
			require.NotEmpty(t, file.Path)
			require.True(t, u.FileExists(file.Path))
		}
	}

	/*
		// try to gc with fewer files than the limit
		{
			err := manager.gc()
			require.NoError(t, err)
		}
	*/

	// create 10 new files
	{
		for i := 0; i < 10; i++ {
			writer, err := manager.New(fmt.Sprintf("hello-%d", i))
			require.NoError(t, err)
			err = writer.Close()
			require.NoError(t, err)
		}
	}

	// check loading files from the directory, should have twelve now
	{
		files, err := manager.Files()
		require.NoError(t, err)
		require.Len(t, files, 10)
		for _, file := range files {
			require.Equal(t, filepath.Dir(file.Path), tempdir)
			require.NotEmpty(t, file.Name)
			require.NotEmpty(t, file.Path)
			require.True(t, u.FileExists(file.Path))
		}
	}

	/*
		// try to gc with fewer files than the limit
		{
			err := manager.gc()
			require.NoError(t, err)
		}
	*/

	// check loading files from the directory, should have ten now
	{
		files, err := manager.Files()
		require.NoError(t, err)
		require.Len(t, files, 10)
	}

	/*
		// try to gc with the current amount of files
		{
			err := logman.LogfileGC(tempdir, 10)
			require.NoError(t, err)
		}

		// check loading files from the directory, should still have ten
		{
			files, err := manager.Files()
			require.NoError(t, err)
			require.Len(t, files, 10)
		}

		// try to gc with only one
		{
			err := logman.LogfileGC(tempdir, 1)
			require.NoError(t, err)
		}

		// check loading files from the directory, should now have only one
		{
			files, err := manager.Files()
			require.NoError(t, err)
			require.Len(t, files, 1)
		}
	*/
}
