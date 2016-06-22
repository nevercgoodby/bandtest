package main

var MsgInterface = []interface{}{
	uint32(1),
	uint16(2),
}

type HeartbeatRequest struct {
	Length  uint32     //len = 12+ 4 + 32 + 64 + 64  == 176
	Cmd     uint32     //(0x00000024), //cmd 80000024 bandtest 0X00000040
	Seq     uint32     //seq
	AppBits uint32     //(0x00000040), //apptype #define OPT_P2PT     0X00000040
	AppName [32]byte   //appname [32]uint8,
	Loads   [32]uint16 //app load on bit 7
	Ports   [32]uint16 //app port on bit 7
}

type HeartbeatResponse struct {
	Length uint32 //len = 12+ 4 + 64 + 64
	Cmd    uint32 //(0x00000024), //cmd 80000024 bandtest 0X00000040
	Seq    uint32 //seq
	Status uint8  //0:ok ,1:error
}
