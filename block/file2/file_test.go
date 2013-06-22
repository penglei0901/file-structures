package file2

import "testing"

import (
    "os"
)

import (
    buf "../buffers"
)

const PATH = "/tmp/__x"

func cleanup(path string) {
    os.Remove(path)
}

func TestOpen(t *testing.T) {
    f := NewBlockFile(PATH, &buf.NoBuffer{})
    defer cleanup(f.Path())
    if err := f.Open(); err != nil {
        t.Fatal(err)
    }
}

func TestAllocate(t *testing.T) {
    f := NewBlockFile(PATH, &buf.NoBuffer{})
    defer cleanup(f.Path())
    if err := f.Open(); err != nil {
        t.Fatal(err)
    }
    if p, err := f.Allocate(); err != nil {
        t.Fatal(err)
    } else if p != 4096 {
        t.Fatalf("Expected p == 4096 got %d", p)
    }
}

func TestSize(t *testing.T) {
    f := NewBlockFile(PATH, &buf.NoBuffer{})
    defer cleanup(f.Path())
    if err := f.Open(); err != nil {
        t.Fatal(err)
    }
    if p, err := f.Allocate(); err != nil {
        t.Fatal(err)
    } else if p != 4096 {
        t.Fatalf("Expected p == 4096 got %d", p)
    }
    if size, err := f.Size(); err != nil {
        t.Fatal(err)
    } else if size != 2*uint64(f.BlkSize()) {
        t.Fatalf("Expected size == %d got %d", 2*f.BlkSize(), size)
    }
}

func TestWriteRead(t *testing.T) {
    f := NewBlockFile(PATH, &buf.NoBuffer{})
    defer cleanup(f.Path())
    if err := f.Open(); err != nil {
        t.Fatal(err)
    }
    if p, err := f.Allocate(); err != nil {
        t.Fatal(err)
    } else if p != 4096 {
        t.Fatalf("Expected p == 4096 got %d", p)
    }
    if size, err := f.Size(); err != nil {
        t.Fatal(err)
    } else if size != 2*uint64(f.BlkSize()) {
        t.Fatalf("Expected size == %d got %d", 2*f.BlkSize(), size)
    }
    blk := make([]byte, f.BlkSize())
    for i := range blk {
        blk[i] = 0xf
    }
    if err := f.WriteBlock(4096, blk); err != nil {
        t.Fatal(err)
    }
    if err := f.Close(); err != nil {
        t.Fatal(err)
    }
    if err := f.Open(); err != nil {
        t.Fatal(err)
    }
    if rblk, err := f.ReadBlock(4096); err != nil {
        t.Fatal(err)
    } else if len(rblk) != int(f.BlkSize()) {
        t.Fatalf("Expected len(rblk) == %d got %d", f.BlkSize(), len(rblk))
    } else {
        for i, b := range rblk {
            if b != 0xf {
                t.Fatalf("Expected rblk[%d] == 0xf got %d", i, b)
            }
        }
    }

    if p, err := f.Allocate(); err != nil {
        t.Fatal(err)
    } else if p != 8192 {
        t.Fatalf("Expected p == 8192 got %d", p)
    }

    if err := f.Free(4096); err != nil {
        t.Fatal(err)
    }
    if p, err := f.Allocate(); err != nil {
        t.Fatal(err)
    } else if p != 4096 {
        t.Fatalf("Expected p == 4096 got %d", p)
    }
    if size, err := f.Size(); err != nil {
        t.Fatal(err)
    } else if size != 3*uint64(f.BlkSize()) {
        t.Fatalf("Expected size == %d got %d", 3*f.BlkSize(), size)
    }
    if err := f.WriteBlock(4096, blk); err != nil {
        t.Fatal(err)
    }
    if err := f.Close(); err != nil {
        t.Fatal(err)
    }
    if err := f.Open(); err != nil {
        t.Fatal(err)
    }
    if rblk, err := f.ReadBlock(4096); err != nil {
        t.Fatal(err)
    } else if len(rblk) != int(f.BlkSize()) {
        t.Fatalf("Expected len(rblk) == %d got %d", f.BlkSize(), len(rblk))
    } else {
        for i, b := range rblk {
            if b != 0xf {
                t.Fatalf("Expected rblk[%d] == 0xf got %d", i, b)
            }
        }
    }
}