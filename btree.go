package btree

import "fmt"

type Btree struct {
	isleaf bool
	Childs []*Btree
	Values []int
	degree int
}

func buildBTreeNode(degree int) *Btree {
	btree := &Btree{
		isleaf: false,
		Childs: []*Btree{},
		Values: []int{},
		degree: degree,
	}
	return btree
}

func buildBtreeRoot(degree int) *Btree {
	btree := &Btree{
		isleaf: true,
		Childs: []*Btree{},
		Values: []int{},
		degree: degree,
	}
	return btree
}

func splitFullChildNode(btree *Btree, i int) {
	// fmt.Println("start")

	// defer func() {
	// 	fmt.Println("end")
	// 	fmt.Printf("%+v", btree)
	// }()
	if len(btree.Childs[i].Values) != btree.MaxNumberOfValues() {
		return
	}
	upVal := btree.Childs[i].Values[btree.degree-1] //上升的值
	newChildNode := buildBTreeNode(btree.degree)
	newChildNode.Values = append(newChildNode.Values, btree.Childs[i].Values[btree.degree:]...)
	newChildNode.isleaf = btree.Childs[i].isleaf
	if !btree.Childs[i].isleaf {
		newChildNode.Childs = append(newChildNode.Childs, btree.Childs[i].Childs[btree.degree:]...)
		btree.Childs[i].Childs = btree.Childs[i].Childs[:btree.degree]
	}
	btree.Childs[i].Values = btree.Childs[i].Values[:btree.degree-1]
	btree.Childs = append(btree.Childs, nil)
	for k := len(btree.Childs) - 2; k > i; k-- {
		btree.Childs[k+1] = btree.Childs[k]
	}
	btree.Childs[i+1] = newChildNode
	btree.Values = append(btree.Values, 0)
	for k := len(btree.Values) - 2; k >= i; k-- {
		btree.Values[k+1] = btree.Values[k]
	}
	btree.Values[i] = upVal
}

func (b *Btree) GetDegree() int {
	return b.degree
}
func (b *Btree) MaxNumberOfChileds() int {
	return b.degree * 2
}

func (b *Btree) MaxNumberOfValues() int {
	return b.degree*2 - 1
}

func (b *Btree) IsFull() bool {
	if b.MaxNumberOfValues() == len(b.Values) {
		return true
	}
	return false
}

func InsertOne(root *Btree, val int) {
	if root.IsFull() {
		newNode := buildBTreeNode(root.degree)
		newNode.isleaf = root.isleaf
		newNode.Childs = root.Childs
		newNode.Values = root.Values
		root.isleaf = false
		root.Childs = []*Btree{newNode}
		root.Values = []int{}
		splitFullChildNode(root, 0)
	}
	recurSplitAndInertOne(root, val)

}

func recurSplitAndInertOne(node *Btree, val int) {
	if node.isleaf {
		insertLeafNode(node, val)
	} else {
		position := 0
		for i := len(node.Values) - 1; i > -1; i-- {
			if node.Values[i] < val {
				position = i + 1
				break
			}
		}
		if node.Childs[position].IsFull() {
			splitFullChildNode(node, position)
			if val > node.Values[position] {
				recurSplitAndInertOne(node.Childs[position+1], val)
			} else {
				recurSplitAndInertOne(node.Childs[position], val)
			}
		} else {
			recurSplitAndInertOne(node.Childs[position], val)
		}
	}
}

func insertLeafNode(node *Btree, val int) {
	node.Values = append(node.Values, 0)
	position := 0
	for i := len(node.Values) - 2; i > -1; i-- {
		if node.Values[i] > val {
			node.Values[i+1] = node.Values[i]
		} else {
			position = i + 1
			break
		}
	}
	node.Values[position] = val
}

//delete
func DeleteOne(root *Btree, val int) {
	//1. judge val in root
	if !root.isleaf { //not leaf node
		if len(root.Values) == 1 && len(root.Childs[0].Values) == root.degree-1 && len(root.Childs[1].Values) == root.degree-1 {
			//merge and heiht-1
			fmt.Println("come here")
			rootVals := root.Childs[0].Values
			rootVals = append(rootVals, root.Values[0])
			rootVals = append(rootVals, root.Childs[1].Values...)
			root.Values = rootVals
			root.isleaf = root.Childs[0].isleaf
			root.Childs = append(root.Childs, root.Childs[1].Childs...)
			root.Childs = root.Childs[0].Childs
		}
	}
	recurMergeAndDeleteOne(root, val)
}

func recurMergeAndDeleteOne(node *Btree, val int) {
	if node.isleaf {
		deleteLeafVal(node, val)
		return
	}
	//inner node
	inNode, posi := isInNode(node, val)
	fmt.Println(inNode, posi)
	if inNode {
		if len(node.Childs[posi].Values) >= node.degree { //先驱
			node.Values[posi] = node.Childs[posi].Values[len(node.Childs[posi].Values)-1]
			node.Childs[posi].Values[len(node.Childs[posi].Values)-1] = val
			recurMergeAndDeleteOne(node.Childs[posi], val)
			return
		} else if len(node.Childs[posi+1].Values) >= node.degree {
			node.Values[posi] = node.Childs[posi+1].Values[0]
			node.Childs[posi+1].Values[0] = val
			recurMergeAndDeleteOne(node.Childs[posi+1], val)
			return
		} else { //merge
			newValue := node.Values[:posi]
			if posi < len(node.Values)-1 {
				newValue = append(newValue, node.Values[posi+1:]...)
			}
			node.Values = newValue
			node.Childs[posi].Values = append(node.Childs[posi].Values, node.Childs[posi+1].Values...)
			if !node.Childs[posi].isleaf {
				node.Childs[posi].Childs = append(node.Childs[posi].Childs, node.Childs[posi+1].Childs...)
			}
			for i := posi + 1; i < len(node.Childs)-1; i++ {
				node.Childs[i] = node.Childs[i+1]
			}
			node.Childs = node.Childs[:len(node.Childs)-1]
			return
		}
	} else {
		fmt.Println(posi)
		if len(node.Childs[posi].Childs) == node.degree {
			if posi+1 < len(node.Childs) && len(node.Childs[posi+1].Childs) > node.degree { //borrow val
				node.Childs[posi].Values = append(node.Childs[posi].Values, node.Values[posi])
				node.Childs[posi].Childs = append(node.Childs[posi].Childs, node.Childs[posi+1].Childs[0])
				node.Values[posi] = node.Childs[posi+1].Values[0]
				node.Childs[posi+1].Values = node.Childs[posi+1].Values[1:]
				node.Childs[posi+1].Childs = node.Childs[posi+1].Childs[1:]
				recurMergeAndDeleteOne(node.Childs[posi], val)
				return
			} else if posi-1 > 0 && len(node.Childs[posi-1].Childs) > node.degree {
				node.Childs[posi].Values = append(node.Childs[posi].Values, node.Values[posi-1])
				node.Childs[posi].Childs = append(node.Childs[posi].Childs, node.Childs[posi-1].Childs[len(node.Childs[posi-1].Childs)-1])
				node.Values[posi] = node.Childs[posi-1].Values[len(node.Childs[posi-1].Values)-1]
				node.Childs[posi-1].Values = node.Childs[posi-1].Values[:len(node.Childs[posi-1].Values)-1]
				node.Childs[posi+1].Childs = node.Childs[posi+1].Childs[:len(node.Childs[posi+1].Childs)-1]
				recurMergeAndDeleteOne(node.Childs[posi], val)
				return
			} else {
				//merge
				if posi+1 >= len(node.Childs) {
					posi -= 1
				}
				newChildValues := node.Childs[posi].Values
				newChildValues = append(newChildValues, node.Values[posi])
				newChildValues = append(newChildValues, node.Childs[posi+1].Values...)
				node.Childs[posi].Childs = append(node.Childs[posi].Childs, node.Childs[posi+1].Childs...)
				node.Childs[posi].Values = newChildValues
				newValues := node.Values[:posi]
				if posi+1 < len(node.Values) {
					newValues = append(newValues, node.Values[posi+1:]...)
				}
				for i := posi + 1; i < len(node.Childs)-1; i++ {
					node.Childs[i] = node.Childs[i+1]
				}
				node.Childs = node.Childs[:len(node.Childs)-1]
				recurMergeAndDeleteOne(node.Childs[posi], val)
				return
			}
		} else {
			recurMergeAndDeleteOne(node.Childs[posi], val)
			return
		}
	}
}

func deleteLeafVal(node *Btree, val int) {
	offset := false
	for i, v := range node.Values {
		if offset {
			node.Values[i-1] = v
			continue
		}
		if v == val {
			offset = true
		}
	}
	node.Values = node.Values[:len(node.Values)-1]
}

func isInNode(node *Btree, val int) (bool, int) {
	position := 0
	inNode := false
	for i := 0; i < len(node.Values); i++ {
		if node.Values[i] < val {
			position++
			continue
		} else if node.Values[i] == val {
			inNode = true
			position = i
			break
		} else {
			inNode = false
			break
		}
	}
	return inNode, position
}
