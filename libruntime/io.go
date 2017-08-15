package libruntime

import (
	"github.com/containerd/containerd"
)

type IO interface {
	Cancel()
	Wait()
	Close()
}

var (
	//IOCreation is just an export
	IOCreation containerd.IOCreation
	NewIO      = containerd.NewIO
)

/*
type IOAttach func(*FIFOSet) (IO, error)


func NewIO(stdin io.Reader, stdout, stderr io.Writer) IOCreation {
	return NewIOWithTerminal(stdin, stdout, stderr, false)
}


func NewIOWithTerminal(stdin io.Reader, stdout, stderr io.Writer, terminal bool) IOCreation {
	return func(id string) (IO, error) {
		paths, err := NewFifos(id)
		if err != nil {
			return nil, err
		}
		i := &IO{
			Terminal: terminal,
			Stdout:   paths.Out,
			Stderr:   paths.Err,
			Stdin:    paths.In,
		}
		set := &ioSet{
			in:  stdin,
			out: stdout,
			err: stderr,
		}
		closer, err := copyIO(paths, set, i.Terminal)
		if err != nil {
			return nil, err
		}
		i.closer = closer
		return i, nil
	}
}

func WithAttach(stdin io.Reader, stdout, stderr io.Writer) IOAttach {
	return func(paths *FIFOSet) (IO, error) {
		if paths == nil {
			return nil, fmt.Errorf("cannot attach to existing fifos")
		}
		i := &IO{
			Terminal: paths.Terminal,
			Stdout:   paths.Out,
			Stderr:   paths.Err,
			Stdin:    paths.In,
		}
		set := &ioSet{
			in:  stdin,
			out: stdout,
			err: stderr,
		}
		closer, err := copyIO(paths, set, i.Terminal)
		if err != nil {
			return nil, err
		}
		i.closer = closer
		return i, nil
	}
}

// Stdio returns an IO implementation to be used for a task
// that outputs the container's IO as the current processes Stdio
func Stdio(id string) (IO, error) {
	return NewIO(os.Stdin, os.Stdout, os.Stderr)(id)
}

// StdioTerminal will setup the IO for the task to use a terminal
func StdioTerminal(id string) (IO, error) {
	return NewIOWithTerminal(os.Stdin, os.Stdout, os.Stderr, true)(id)
}

type FIFOSet struct {
	// Dir is the directory holding the task fifos
	Dir          string
	In, Out, Err string
	Terminal     bool
}

type ioSet struct {
	in       io.Reader
	out, err io.Writer
}

type wgCloser struct {
	wg     *sync.WaitGroup
	dir    string
	set    []io.Closer
	cancel context.CancelFunc
}

func (g *wgCloser) Wait() {
	g.wg.Wait()
}

func (g *wgCloser) Close() error {
	for _, f := range g.set {
		f.Close()
	}
	if g.dir != "" {
		return os.RemoveAll(g.dir)
	}
	return nil
}

func (g *wgCloser) Cancel() {
	g.cancel()
}

// NewFifos returns a new set of fifos for the task
func NewFifos(id string) (*FIFOSet, error) {
	root := filepath.Join(os.TempDir(), "containerd")
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, err
	}
	dir, err := ioutil.TempDir(root, "")
	if err != nil {
		return nil, err
	}
	return &FIFOSet{
		Dir: dir,
		In:  filepath.Join(dir, id+"-stdin"),
		Out: filepath.Join(dir, id+"-stdout"),
		Err: filepath.Join(dir, id+"-stderr"),
	}, nil
}

func copyIO(fifos *FIFOSet, ioset *ioSet, tty bool) (_ *wgCloser, err error) {
	var (
		f           io.ReadWriteCloser
		set         []io.Closer
		ctx, cancel = context.WithCancel(context.Background())
		wg          = &sync.WaitGroup{}
	)
	defer func() {
		if err != nil {
			for _, f := range set {
				f.Close()
			}
			cancel()
		}
	}()

	if f, err = fifo.OpenFifo(ctx, fifos.In, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
		return nil, err
	}
	set = append(set, f)
	go func(w io.WriteCloser) {
		io.Copy(w, ioset.in)
		w.Close()
	}(f)

	if f, err = fifo.OpenFifo(ctx, fifos.Out, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
		return nil, err
	}
	set = append(set, f)
	wg.Add(1)
	go func(r io.ReadCloser) {
		io.Copy(ioset.out, r)
		r.Close()
		wg.Done()
	}(f)

	if f, err = fifo.OpenFifo(ctx, fifos.Err, syscall.O_RDONLY|syscall.O_CREAT|syscall.O_NONBLOCK, 0700); err != nil {
		return nil, err
	}
	set = append(set, f)

	if !tty {
		wg.Add(1)
		go func(r io.ReadCloser) {
			io.Copy(ioset.err, r)
			r.Close()
			wg.Done()
		}(f)
	}
	return &wgCloser{
		wg:     wg,
		dir:    fifos.Dir,
		set:    set,
		cancel: cancel,
	}, nil
}
*/
