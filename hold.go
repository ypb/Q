// Copyright 2011 AUTHORS. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package main

// Q type is
type Q interface {
	Is() interface{} // identity
	Hold(Q)          // b(contain)/p(contains)/r(...
	Heldby(Q) Q      // p(contained)/b(containing)/r(...
	Print()
}
// p type exists
type p struct {
	id      interface{}
	holding *b
	heldby  *p
}
// b type binds?
type b struct {
	id      *p
	holding Q
	heldby  *b
}
// r type reflects? relates...? WTF, it's only a linked list...
type r struct {
	id      *b
	holding *p
	heldby  Q
}

func (q p) Is() interface{} {
	return q.id
}
func (d b) Is() interface{} {
	return interface{}(d.holding)
}
func (s r) Is() interface{} {
	return interface{}(s.heldby)
}

func (q *p) Hold(ob Q) {
	if q.holding == nil {
		q.holding = &b{q, ob, nil}
		return
	}
	q.holding.Hold(ob)
}
func (d *b) Hold(ob Q) {
	if d.heldby == nil {
		d.heldby = &b{d.id, ob, nil}
		return
	}
	d.heldby.Hold(ob)
}
func (s r) Hold(ob Q) {
	// place holder
	return
}

func (q *p) Heldby(ob Q) Q {
	if q.heldby == nil {
		q.heldby = ob.(*p)
		return q
	}
	if q.holding == nil {
		// return nil
		// idiot?!?
		q.holding = &b{q, nil, nil}
	}
	return q.holding.Heldby(ob)
}
func (d *b) Heldby(ob Q) Q {
	if d.holding == nil {
		d.holding = r{d, nil, ob}
		return d.id
	}
	if d.holding.(Q).Is() == ob.Is() {
		return d.holding
	}
	if d.heldby == nil {
		// do the 2. p.Heldby and 2. b.Heldby HERE?
		return nil
	}
	return d.heldby.Heldby(ob)
}
func (s r) Heldby(ob Q) Q {
	// place holder
	return ob
}

func (q p) Print() {
	end := q.holding
	str := " " + q.id.(string)
	for end != nil {
		println(str, end.holding.Is().(string))
		end = end.heldby
	}
}
func (d b) Print() {
	println(d.holding)
}
func (s r) Print() {
	println(s.holding)
}


func P(id interface{}) Q {
	return &p{id, nil, nil}
}

// Local Variables:
// mode: Go
// End:
