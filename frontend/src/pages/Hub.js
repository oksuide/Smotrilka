import React, { useState } from 'react';
import { FaCog, FaPlus, FaSearch } from 'react-icons/fa';
import { Link } from 'react-router-dom';

function HubPage() {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [roomName, setRoomName] = useState('');
    const [roomPassword, setRoomPassword] = useState('');

    // Sample rooms data
    const [rooms, setRooms] = useState([]);

    const handleCreateRoomClick = () => {
        setIsModalOpen(true);
    };

    const handleModalClose = () => {
        setIsModalOpen(false);
        setRoomName('');
        setRoomPassword('');
    };

    const handleCreateRoom = () => {
        console.log('Room Name:', roomName);
        console.log('Room Password:', roomPassword);
        handleModalClose();
    };

    return (
        <div style={backgroundStyle}>
            <div style={leftBlockStyle}>
                <div style={userContainerStyle}>
                    <div style={avatarStyle}></div>
                    <Link to="/profile" style={usernameStyle}>Username</Link>
                </div>
                <h2 style={myRoomsStyle}>My rooms</h2>
                <div style={searchContainerStyle}>
                    <input type="text" placeholder="Search room" style={searchInputStyle} />
                    <FaSearch style={searchIconStyle} />
                </div>
                <div style={roomListStyle}>
                    {rooms.map((room, index) => (
                        <RoomItem key={index} name={room.name} participants={room.participants} />
                    ))}
                    <div style={createRoomStyle} onClick={handleCreateRoomClick}>
                        <FaPlus style={plusIconStyle} />
                        <span>Create new room</span>
                    </div>
                </div>
            </div>
            <div style={rightBlockStyle}>
                <FaCog style={settingsIconStyle} />
                <h2 style={joinRoomStyle}>Join to room</h2>
                <input type="text" placeholder="Room ID" style={smallInputStyle} />
                <input type="password" placeholder="Password (if needed)" style={smallInputStyle} />
                <button style={connectButtonStyle}>Connect</button>
            </div>
            {isModalOpen && (
                <div style={modalOverlayStyle}>
                    <div style={modalStyle}>
                        <h2>Enter the room name and password</h2>
                        <input
                            type="text"
                            placeholder="Room Name"
                            value={roomName}
                            onChange={(e) => setRoomName(e.target.value)}
                            style={modalInputStyle}
                        />
                        <input
                            type="password"
                            placeholder="Password"
                            value={roomPassword}
                            onChange={(e) => setRoomPassword(e.target.value)}
                            style={modalInputStyle}
                        />
                        <button onClick={handleCreateRoom} style={createButtonStyle}>Create</button>
                        <button onClick={handleModalClose} style={cancelButtonStyle}>Cancel</button>
                    </div>
                </div>
            )}
        </div>
    );
}

function RoomItem({ name, participants }) {
    return (
        <div style={roomItemStyle}>
            <div style={roomAvatarStyle}></div>
            <span style={roomNameStyle}>{name}</span>
            <span>{participants}</span>
        </div>
    );
}

const backgroundStyle = {
    display: 'flex',
    height: '100vh',
    backgroundColor: '#1e1e1e',
    color: 'white',
    overflow: 'hidden',
};

const leftBlockStyle = {
    width: '50%',
    padding: '30px',
    borderRight: '1.5px solid white',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'flex-start',
};

const rightBlockStyle = {
    width: '50%',
    padding: '30px',
    position: 'relative',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
};

const userContainerStyle = {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '30px',
};

const avatarStyle = {
    width: '60px',
    height: '60px',
    borderRadius: '50%',
    backgroundColor: 'gray',
    marginRight: '15px',
};

const usernameStyle = {
    fontSize: '27px',
    color: 'white',
    textDecoration: 'none',
};

const myRoomsStyle = {
    color: 'white',
    fontSize: '27px',
};

const searchContainerStyle = {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '30px',
};

const searchInputStyle = {
    width: '40%',
    padding: '12px',
    borderRadius: '7.5px',
    border: '1.5px solid white',
    backgroundColor: '#333',
    color: 'white',
};

const searchIconStyle = {
    marginLeft: '15px',
};

const roomListStyle = {
    display: 'flex',
    flexDirection: 'column',
};

const roomItemStyle = {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    padding: '15px',
    borderBottom: '1.5px solid white',
};

const roomAvatarStyle = {
    width: '60px',
    height: '60px',
    backgroundColor: 'gray',
    marginRight: '15px',
};

const roomNameStyle = {
    flex: 1,
    marginLeft: '15px',
};

const createRoomStyle = {
    display: 'flex',
    alignItems: 'center',
    padding: '15px',
    marginTop: '30px',
    borderTop: '1.5px solid white',
    cursor: 'pointer',
};

const plusIconStyle = {
    color: 'red',
    marginRight: '15px',
};

const joinRoomStyle = {
    marginBottom: '30px',
    fontSize: '27px',
};

const smallInputStyle = {
    width: '60%',
    padding: '12px',
    marginBottom: '15px',
    borderRadius: '7.5px',
    border: '1.5px solid white',
    backgroundColor: '#333',
    color: 'white',
};

const connectButtonStyle = {
    padding: '15px 30px',
    borderRadius: '7.5px',
    border: 'none',
    backgroundColor: '#337a3c',
    color: 'white',
    cursor: 'pointer',
};

const settingsIconStyle = {
    position: 'absolute',
    top: '30px',
    right: '30px',
    cursor: 'pointer',
};

const modalOverlayStyle = {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
};

const modalStyle = {
    backgroundColor: '#1e1e1e',
    padding: '20px',
    borderRadius: '10px',
    width: '280px', // Reduced width
    textAlign: 'center',
    border: '2px solid #337a3c', // Dark purple border
};

const modalInputStyle = {
    width: '90%', // Reduced width
    padding: '8px', // Reduced padding
    margin: '10px 0',
    borderRadius: '5px',
    border: '1px solid #ccc',
};

const createButtonStyle = {
    padding: '10px 20px',
    margin: '10px 5px',
    borderRadius: '5px',
    border: 'none',
    backgroundColor: '#138513', // Pastel green
    color: 'white',
    cursor: 'pointer',
};

const cancelButtonStyle = {
    padding: '10px 20px',
    margin: '10px 5px',
    borderRadius: '5px',
    border: 'none',
    backgroundColor: '#8c0b0b', // Pastel red
    color: 'white',
    cursor: 'pointer',
};

export default HubPage;