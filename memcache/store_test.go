package memcache

import (
	"bytes"
	"fmt"
	"testing"
)

func testDistributeStore(t *testing.T, dclient DistributeStorage) {
	key := "dtest"
	r, hs, _ := dclient.Get(key)
	if r != nil && len(hs) > 0 {
		t.Errorf("get should return nil and []")
	}
	// set
	v := []byte("value")
	flag := 2
	ok, hs, _ := dclient.Set(key, &Item{Body: v, Flag: flag}, false)
	if !ok || len(hs) == 0 {
		t.Errorf("set failed")
	}
	v2, h, _ := dclient.Get(key)
	if v2 == nil || !bytes.Equal(v, v2.Body) || len(h) == 0 {
		t.Errorf("should return same value and a get host")
	} else {
		in_host := false
		for _, set_h := range hs {
			if set_h == h {
				in_host = true
				break
			}
		}
		if !in_host {
			t.Errorf("get should from set hosts")
		}
	}
	if v2.Flag != flag {
		t.Errorf("should return flag 2")
	}
	// set with noreply
	v = []byte("value 2")
	flag = 3
	key2 := "test2"
	ok, hs, _ = dclient.Set(key2, &Item{Body: v, Flag: flag}, true)
	if !ok || len(hs) == 0 {
		t.Errorf("set with ply failed")
	}
	v2, h, _ = dclient.Get(key2)
	if v2 == nil || !bytes.Equal(v, v2.Body) || len(h) == 0 {
		t.Errorf("should return same value")
	} else {
		in_host := false
		for _, set_h := range hs {
			if set_h == h {
				in_host = true
				break
			}
		}
		if !in_host {
			t.Errorf("get should from set hosts")
		}
	}
	// get_multi
	items, hhs, _ := dclient.GetMulti([]string{"test", "test", "test2", "test3"})
	if len(items) != 2 || len(hhs) == 0 {
		t.Errorf("get_multi should return 2 values, but got %d", len(items))
	}
	keys := make([]string, 102)
	for i := 0; i < 100; i++ {
		keys[i] = fmt.Sprintf("__t%d", i)
		dclient.Set(keys[i], &Item{Body: v}, true)
	}
	items, hhs, _ = dclient.GetMulti(keys)
	if len(items) != 100 || len(hhs) == 0 {
		t.Errorf("get_multi should return 100 values, but got %d", len(items))
	}
	// get large obj
	v = make([]byte, 1024*1000)
	if ok, hhs, _ := dclient.Set("test_large", &Item{Body: v, Flag: flag}, false); !ok || len(hhs) == 0 {
		t.Errorf("set large value failed")
	}
	v2, _ = dclient.Get("test_large")
	if v2 == nil || !bytes.Equal(v, v2.Body) {
		t.Errorf("should return same large value")
	}
	// append
}

func testStore(t *testing.T, client Storage) {
	key := "test"
	r, _ := client.Get("test")
	if r != nil {
		t.Errorf("get should return nil")
	}
	// set
	v := []byte("value")
	flag := 2
	if ok, _ := client.Set("test", &Item{Body: v, Flag: flag}, false); !ok {
		t.Errorf("set failed")
	}
	v2, _ := client.Get("test")
	if !bytes.Equal(v, v2.Body) {
		t.Errorf("should return value")
	}
	if v2.Flag != flag {
		t.Errorf("should return flag 2")
	}
	// set with noreply
	v = []byte("value 2")
	flag = 3
	if ok, _ := client.Set("test2", &Item{Body: v, Flag: flag}, true); !ok {
		t.Errorf("set failed")
	}
	v2, _ = client.Get("test2")
	if v2 == nil || !bytes.Equal(v, v2.Body) {
		t.Errorf("should return value")
	}
	// get_multi
	items, _ := client.GetMulti([]string{"test", "test", "test2", "test3"})
	if len(items) != 2 {
		t.Errorf("get_multi should return 2 values, but got %d", len(items))
	}
	keys := make([]string, 102)
	for i := 0; i < 100; i++ {
		keys[i] = fmt.Sprintf("__t%d", i)
		client.Set(keys[i], &Item{Body: v}, true)
	}
	items, _ = client.GetMulti(keys)
	if len(items) != 100 {
		t.Errorf("get_multi should return 100 values, but got %d", len(items))
	}
	// get large obj
	v = make([]byte, 1024*1000)
	if ok, _ := client.Set("test_large", &Item{Body: v, Flag: flag}, false); !ok {
		t.Errorf("set large value failed")
	}
	v2, _ = client.Get("test_large")
	if v2 == nil || !bytes.Equal(v, v2.Body) {
		t.Errorf("should return large value")
	}
	// append
	client.Set(key, &Item{Body: []byte("value")}, false)
	if ok, _ := client.Append("test", []byte(" good")); !ok {
		t.Error("append failed")
	}
	v2, _ = client.Get("test")
	if v2 == nil || string(v2.Body) != "value good" {
		t.Errorf("get after append: %v", v2)
	}
	// incr
	client.Set("test", &Item{Body: []byte("3"), Flag: 4}, false)
	if v, _ := client.Incr("test", 5); v != 8 {
		t.Errorf("incr failed: %d!=8", v)
	}
	// delete
	if ok, _ := client.Delete("test"); !ok {
		t.Errorf("delete failed")
	}
	v2, _ = client.Get("test")
	if v2 != nil {
		t.Errorf("get should return []")
	}
}

func testFailStore(t *testing.T, store Storage) {
	_, err := store.Get("key")
	if err == nil {
		t.Error("Get() should raise error")
	}
	_, err = store.GetMulti([]string{"key"})
	if err == nil {
		t.Error("GetMulti() should raise error")
	}
	_, err = store.Set("key", &Item{}, false)
	if err == nil {
		t.Error("Set() should raise error")
	}
	_, err = store.Append("key", nil)
	if err == nil {
		t.Error("Append() should raise error")
	}
	_, err = store.Incr("key", 1)
	if err == nil {
		t.Error("Incr() should raise error")
	}
	_, err = store.Delete("key")
	if err == nil {
		t.Error("Delete() should raise error")
	}
}

func TestStore(t *testing.T) {
	store := NewMapStore()
	testStore(t, store)
}
