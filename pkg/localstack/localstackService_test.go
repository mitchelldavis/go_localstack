package localstack

import (
    "fmt"
    "log"
    "testing"
)

func Test_NewLocalstackService(t *testing.T) {
	expected := &LocalstackService {
		Name: "sqs",
		Protocol: "tcp",
		Port: 4576,
	}
	lss, err := NewLocalstackService("sqs")

	if err != nil {
		log.Fatal("An error was not expected with a service request of: sqs")
	}

	if	!expected.Equals(lss) {
		log.Fatal("The resulting service did not equal the expected service.")
	}

	lss, err = NewLocalstackService("garbage")

	if err == nil {
		log.Fatal("An error was expected with a service request of: garbage")
	}
}

func Test_Localstack_Equal(t *testing.T) {
	var actual, expected *LocalstackService

	if !actual.Equals(expected) {
		log.Fatal("When both pointers are null, the result should be true.")
	}

	expected = &LocalstackService {
		Name: "sqs",
		Protocol: "tcp",
		Port: 4576,
	}

	if actual.Equals(expected) {
		log.Fatal("When the lhs is nil but the rhs is not, then the result should be false.")
	}

	actual = &LocalstackService {
		Name: "es",
		Protocol: "tcp",
		Port: 4576,
	}

	if actual.Equals(expected) {
		log.Fatal("When the lhs.Name != rhs.Name the result should be false.")
	}

	actual = &LocalstackService {
		Name: "sqs",
		Protocol: "udp",
		Port: 4576,
	}

	if actual.Equals(expected) {
		log.Fatal("When the lhs.Protocol != rhs.Protocol the result should be false.")
	}

	actual = &LocalstackService {
		Name: "sqs",
		Protocol: "tcp",
		Port: 0,
	}

	if actual.Equals(expected) {
		log.Fatal("When the lhs.Port != rhs.Port the result should be false.")
	}

	actual = &LocalstackService {
		Name: "es",
		Protocol: "udp",
		Port: 0,
	}

	if actual.Equals(expected) {
		log.Fatal("When no fields in lhs equal rhs' fields the result should be false.")
	}

	actual = &LocalstackService {
		Name: "sqs",
		Protocol: "tcp",
		Port: 4576,
	}

	if !actual.Equals(expected) {
		log.Fatal("When all the fields in lhs equal rhs' fields the result should be true.")
	}
}

func Test_LocalstackService_GetPortProtocol(t *testing.T) {
	actual := &LocalstackService {
		Name: "sqs",
		Protocol: "tcp",
		Port: 4576,
	}

	if actual.GetPortProtocol() != fmt.Sprintf("%d/%s", actual.Port, actual.Protocol) {
		log.Fatal(fmt.Sprintf("%s != %s", actual.GetPortProtocol(), fmt.Sprintf("%d:%s", actual.Port, actual.Protocol)))
	}
}

func Test_LocalstackService_GetNamePort(t *testing.T) {
	actual := &LocalstackService {
		Name: "sqs",
		Protocol: "tcp",
		Port: 4576,
	}

	if actual.GetNamePort() != fmt.Sprintf("%s:%d", actual.Name, actual.Port) {
		log.Fatal(fmt.Sprintf("%s != %s", actual.GetNamePort(), fmt.Sprintf("%s:%d", actual.Name, actual.Port)))
	}
}

func Test_LocalstackServiceCollection_GetServiceMap(t *testing.T) {
	first, _ := NewLocalstackService("sqs")
	second, _ := NewLocalstackService("sns")
	lsc := LocalstackServiceCollection {
		*first,
		*second,
	}

	expected := fmt.Sprintf("sqs:4576,sns:4575")

	if lsc.GetServiceMap() != expected {
		log.Fatal(fmt.Sprintf("%s != %s", lsc.GetServiceMap(), expected))
	}
}

func Test_LocalstackServiceCollection_Len(t *testing.T) {
	first, _ := NewLocalstackService("sqs")
	second, _ := NewLocalstackService("sns")
	lsc := LocalstackServiceCollection {
		*first,
		*second,
	}

	expected := 2

	if lsc.Len() != expected {
		log.Fatal(fmt.Sprintf("%d != %d", lsc.Len(), expected))
	}

	third, _ := NewLocalstackService("es")
	lsc = LocalstackServiceCollection {
		*first,
		*second,
		*third,
	}

	expected = 3

	if lsc.Len() != expected {
		log.Fatal(fmt.Sprintf("%d != %d", lsc.Len(), expected))
	}
}

func Test_LocalstackServiceCollection_Swap(t *testing.T) {
	first, _ := NewLocalstackService("sqs")
	second, _ := NewLocalstackService("sns")
	third, _ := NewLocalstackService("es")
	lsc := LocalstackServiceCollection {
		*first,
		*second,
		*third,
	}

	lsc.Swap(0, 1)

	if !(lsc[0].Equals(second) && lsc[1].Equals(first)) {
		log.Fatal("The first and second item should have been swapped.")
	}

	lsc.Swap(1, 2)

	if !(lsc[1].Equals(third) && lsc[2].Equals(first)) {
		log.Fatal("The second and third item should have been swapped.")
	}
}

func Test_LocalstackServiceCollection_Less(t *testing.T) {
	first, _ := NewLocalstackService("sqs")
	second, _ := NewLocalstackService("sns")
	third, _ := NewLocalstackService("es")
	lsc := LocalstackServiceCollection {
		*first,
		*second,
		*third,
	}

	if lsc.Less(0, 1) {
		log.Fatal("The first item should be less than the second")
	}

	if lsc.Less(1, 2) {
		log.Fatal("The second item should be less than the third")
	}

	if !lsc.Less(2, 0) {
		log.Fatal("The third item should be less than the first")
	}
}

func Test_LocalstackServiceCollection_Sort(t *testing.T) {
	first, _ := NewLocalstackService("sqs")
	second, _ := NewLocalstackService("sns")
	third, _ := NewLocalstackService("es")
	lsc := LocalstackServiceCollection {
		*first,
		*second,
		*third,
	}

	expected := LocalstackServiceCollection {
		*third,
		*second,
		*first,
	}

	lsc.Sort()

	if !(lsc[0].Equals(&expected[0]) && lsc[1].Equals(&expected[1]) && lsc[2].Equals(&expected[2])) {
		log.Fatal("The sort order isn't matching what is expected.")
	}
}

func Test_LocalstackServiceCollection_Contains(t *testing.T) {
	s3, _ := NewLocalstackService("s3")
	sqs, _ := NewLocalstackService("sqs")

    lsc := LocalstackServiceCollection {
        *s3,
        *sqs,
    }

    if !lsc.Contains("s3") {
        t.Error("s3 was added to the collection but Contains says it was not.")
    }
    if !lsc.Contains("sqs") {
        t.Error("sqs was added to the collection but Contains says it isn't there.")
    }
    if lsc.Contains("redshift") {
        t.Error("redshift was not added to the collection but Contains says it was.")
    }
    if lsc.Contains("dynamodb") {
        t.Error("dynamodb was not added to the collection but Contains says it was.")
    }
}
