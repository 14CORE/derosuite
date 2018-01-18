package p2p

import "encoding/binary"

var FLAGS_VALUE uint32 = 0 // we donot support fluffly blocks at this point in time
var support_response_bytes = [29]byte{0x01, 0x11, 0x01, 0x01, 0x01, 0x01, 0x02, 0x01, 0x01, 0x04, 0x0d, 0x73, 0x75, 0x70, 0x70,
	0x6f, 0x72, 0x74, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x06, 0x00, 0x00, 0x00, 0x00}

// handle P2P_COMMAND_REQUEST_SUPPORT_FLAGS
func Handle_P2P_Support_Flags(connection *Connection,
	i_command_header *Levin_Header, buf []byte) {

	// deserialize data header
	var i_data_header Levin_Data_Header // incoming data header

	err := i_data_header.DeSerialize(buf)

	if err != nil {
		connection.logger.Debugf("We should destroy connection here, data header cnot deserialized")
		connection.Exit = true
		return
	}
	// make sure  data is length 10
	// create a new response header

	var o_command_header Levin_Header
	//var o_data_header  Levin_Data_Header

	binary.LittleEndian.PutUint32(support_response_bytes[25:], FLAGS_VALUE) // packed 8+8+1+4 bytes

	o_command_header.CB = uint64(len(support_response_bytes))

	o_command_header.Command = P2P_COMMAND_REQUEST_SUPPORT_FLAGS
	o_command_header.ReturnData = false
	o_command_header.Flags = LEVIN_PACKET_RESPONSE

	o_command_header_bytes, _ := o_command_header.Serialize()

	connection.Conn.Write(o_command_header_bytes)
	connection.Conn.Write(support_response_bytes[:])

}

// send the hand shake
func Send_SupportFlags_Command(connection *Connection) {

	connection.Lock()

	var o_command_header Levin_Header
	var o_data_header Levin_Data_Header

	o_data_bytes, _ := o_data_header.Serialize()

	o_command_header.CB = uint64(len(o_data_bytes))

	o_command_header.Command = P2P_COMMAND_REQUEST_SUPPORT_FLAGS
	o_command_header.ReturnData = true
	o_command_header.Flags = LEVIN_PACKET_REQUEST

	o_command_header_bytes, _ := o_command_header.Serialize()

	connection.Conn.Write(o_command_header_bytes)
	connection.Conn.Write(o_data_bytes)

	connection.Command_queue.PushBack(uint32(P2P_COMMAND_REQUEST_SUPPORT_FLAGS))

	connection.Unlock()
}
