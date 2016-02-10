package main

import (
	"list"
	"net"
)



const BucketSize = 20;

type Contact struct {
  id NodeID
  ip net.IP

}

type RoutingTable struct {
  node NodeID;
  buckets [IdLength*8]*list.List;
}

func NewRoutingTable(node NodeID) (ret RoutingTable) {
  for i := 0; i < IdLength * 8; i++ {
    ret.buckets[i] = list.New();
  }
  ret.node = node;
  return;
}