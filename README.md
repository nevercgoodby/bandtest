Golang中的字节序列化操作
2014-08-12 20:01:54
标签：字节 golang
原创作品，允许转载，转载时请务必以超链接形式标明文章 原始出处 、作者信息和本声明。否则将追究法律责任。http://gotaly.blog.51cto.com/8861157/1539119
在写网络程序的时候，我们经常需要将结构体或者整数等数据类型序列化成二进制的buffer串。或者从一个buffer中解析出来一个结构体出来，最典型的就是在协议的header部分表征head length 或者body length在拼包和拆包的过程中，需要按照规定的整数类型进行解析，且涉及到大小端序的问题。

1.C中是怎么操作的
在C中我们最简单的方法是用memcpy来一个整形数或者结构体等其他类型复制到一块内存中，然后在强转回需要的类型。如:

    // produce
    int a = 32;
    char *buf  = (char *)malloc(sizeof(int));
    memcpy(buf,&a,sizeof(int));

    // consume
    int b ;
    memcpy(&b,buf,sizeof(int))
必要的时候采用ntoh/hton系列函数进行大小端序的转换。

2.golang中操作
通过"encoding/binary"可以提供常用的二进制序列化的功能。该模块主要提供了如下几个接口：

func Read(r io.Reader, order ByteOrder, data interface{}) error
func Write(w io.Writer, order ByteOrder, data interface{}) error
func Size(v interface{}) int

var BigEndian bigEndian
var LittleEndian littleEndian
/*
type ByteOrder interface {
Uint16([]byte) uint16
Uint32([]byte) uint32
Uint64([]byte) uint64
PutUint16([]byte, uint16)
PutUint32([]byte, uint32)
PutUint64([]byte, uint64)
String() string
}
/*
通过Read接口可以将buf中得内容填充到data参数表示的数据结构中，通过Write接口可以将data参数里面包含的数据写入到buffer中。 变量BigEndian和LittleEndian是实现了ByteOrder接口的对象，通过接口中提供的方法可以直接将uintx类型序列化（uintx()）或者反序列化(putuintx())到buf中。

2.1将结构体序列化到一个buf中
在序列化结构对象时，需要注意的是，被序列化的结构的大小必须是已知的，可以通过Size接口来获得该结构的大小，从而决定buffer的大小。

i := uint16(1)
size :=  binary.Size(i)
固定大小的结构体，就要求结构体中不能出现[]byte这样的切片成员，否则Size返回-1，且不能进行正常的序列化操作。

type A struct {
    // should be exported member when read back from buffer
    One int32
    Two int32
}

var a A


a.One = int32(1)
a.Two = int32(2)

buf := new(bytes.Buffer)
fmt.Println("a's size is ",binary.Size(a))
binary.Write(buf,binary.LittleEndian,a)
fmt.Println("after write ，buf is:",buf.Bytes())
对应的输出为：

a's size is  8
after write ,buf is : [1 0 0 0 2 0 0 0]
通过Size可以得到所需buffer的大小。通过Write可以将对象a的内容序列化到buffer中。这里采用了小端序的方式进行序列化（x86架构都是小端序，网络字节序是大端序）。

对于结构体中得“_”成员不进行序列化。

2.2从buf中反序列化回一个结构
从buffer中读取时，一样要求结构体的大小要固定，且需要反序列化的结构体成员必须是可导出的也就是必须是大写开头的成员，同样对于“_”不进行反序列化：

type A struct {
    // should be exported member when read back from buffer
    One int32
    Two int32
}

var aa A

buf := new(bytes.Buffer)
binary.Write(buf,binary.LittleEndian,a)
binary.Read(buf,binary.LittleEndian,&aa)
fmt.Println("after aa is ",aa)
输出为：

after write ,bufis : [1 0 0 0 2 0 0 0]
before aa is : {0 0}
after aa is  {1 2}
这里使用Read从buffer中将数据导入到结构体对象aa中。如果结构体中对应的成员不是可导出的，那么在转换的时候会panic出错。

2.3将整数序列化到buf中，并从buf中反序列化出来
我们可以通过Read/Write直接去读或者写一个uintx类型的变量来实现对整形数的序列化和反序列化。由于在网络中，对于整形数的序列化非常常用，因此系统库提供了type ByteOrder接口可以方便的对uint16/uint32/uint64进行序列化和反序列化：

int16buf := new(bytes.Buffer)
i := uint16(1)
binary.Write(int16buf,binary.LittleEndian,i)
fmt.Println(“write buf is:”int16buf.Bytes())

var int16buf2 [2]byte
binary.LittleEndian.PutUint16(int16buf2[:],uint16(1))
fmt.Println("put buffer is :",int16buf2[:])

ii := binary.LittleEndian.Uint16(int16buf2[:])
fmt.Println("Get buf is :",ii)
输出为：

write buffer is : [1 0]
put buf is: [1 0]
Get buf is : 1
通过调用binary.LittleEndian.PutUint16,可以按照小端序的格式将uint16类型的数据序列化到buffer中。通过binary.LittleEndian.Uint16将buffer中内容反序列化出来。

3. 一个实在的例子
我们来看一个网络包包头的定义和初始化：

type Head struct {
    Cmd byte
    Version byte
    Magic   uint16
    Reserve byte
    HeadLen byte
    BodyLen uint16
}

func NewHead(buf []byte)*Head{
    head := new(Head)

    head.Cmd     = buf[0]
    head.Version = buf[1]
    head.Magic   = binary.BigEndian.Uint16(buf[2:4])
    head.Reserve = buf[4]
    head.HeadLen = buf[5]
    head.BodyLen = binary.BigEndian.Uint16(buf[6:8])
    return head
}
这个是一个常见的在tcp 拼包得例子。在例子中通过binary.BigEndian.Uint16将数据按照网络序的格式读出来，放入到head中对应的结构里面。


本文出自 “Done_in_72_hours” 博客，请务必保留此出处http://gotaly.blog.51cto.com/8861157/1539119