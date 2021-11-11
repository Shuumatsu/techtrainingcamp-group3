package models

func int2str(num uint64) string {
	if num == 0 {
		return "0"
	}
	var ret []byte
	for num != 0 {
		ret = append(ret, byte(num%10)+'0')
		num /= 10
	}
	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}
	return string(ret)
}
