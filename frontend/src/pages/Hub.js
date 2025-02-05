import axios from 'axios';
import React, { useEffect, useRef, useState } from 'react';
import { FaCog, FaPlus, FaSearch } from 'react-icons/fa';
import { Link } from 'react-router-dom';
import styles from '../styles/HubPageStyles';


function HubPage() {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [roomName, setRoomName] = useState('');
    const [roomPassword, setRoomPassword] = useState('');
    const [rooms, setRooms] = useState([]);
    const [myRooms, setMyRooms] = useState([]);
    const [username, setUsername] = useState('');
    const effectRan = useRef(false);

    useEffect(() => {
        if (effectRan.current) return;
        effectRan.current = true;

        const token = localStorage.getItem('token');
        if (token) {
            axios.get('http://localhost:8080/api/users/me', {
                headers: { 'Authorization': `Bearer ${token}` }
            }).then(response => {
                setUsername(response.data.username);
            }).catch(error => {
                console.error('Error fetching user data:', error);
            });

            axios.get('http://localhost:8080/api/rooms/my', {
                headers: { 'Authorization': `Bearer ${token}` }
            }).then(response => {
                setMyRooms(response.data);
            }).catch(error => {
                console.error('Error fetching user rooms:', error);
            });
        }
    }, []);

    const handleCreateRoomClick = () => {
        setIsModalOpen(true);
    };

    const handleModalClose = () => {
        setIsModalOpen(false);
        setRoomName('');
        setRoomPassword('');
    };

    const handleCreateRoom = () => {
        const token = localStorage.getItem('token');
        axios.post('http://localhost:8080/api/rooms', { name: roomName, password: roomPassword }, {
            headers: { 'Authorization': `Bearer ${token}` }
        })
            .then(response => {
                setMyRooms([...myRooms, response.data]);
                handleModalClose();
            })
            .catch(error => {
                console.error('Error creating room:', error);
            });
    };

    return (
        <div style={styles.backgroundStyle}>
            <div style={styles.leftBlockStyle}>
                <div style={styles.userContainerStyle}>
                    <div style={styles.avatarStyle}></div>
                    <Link to="/profile/:id" style={styles.usernameStyle}>{username || 'Loading...'}</Link>
                </div>
                <h2 style={styles.myRoomsStyle}>My rooms</h2>
                <div style={styles.roomListStyle}>
                    {myRooms.map((room, index) => (
                        <RoomItem key={index} name={room.name} participants={room.participants} />
                    ))}
                </div>
                <h2 style={styles.myRoomsStyle}>All rooms</h2>
                <div style={styles.searchContainerStyle}>
                    <input type="text" placeholder="Search room" style={styles.searchInputStyle} />
                    <FaSearch style={styles.searchIconStyle} />
                </div>
                <div style={styles.roomListStyle}>
                    {rooms.map((room, index) => (
                        <RoomItem key={index} name={room.name} participants={room.participants} />
                    ))}
                    <div style={styles.createRoomStyle} onClick={handleCreateRoomClick}>
                        <FaPlus style={styles.plusIconStyle} />
                        <span>Create new room</span>
                    </div>
                </div>
            </div>
            <div style={styles.rightBlockStyle}>
                <FaCog style={styles.settingsIconStyle} />
                <h2 style={styles.joinRoomStyle}>Join to room</h2>
                <input type="text" placeholder="Room ID" style={styles.smallInputStyle} />
                <input type="password" placeholder="Password (if needed)" style={styles.smallInputStyle} />
                <button style={styles.connectButtonStyle}>Connect</button>
            </div>
            {isModalOpen && (
                <div style={styles.modalOverlayStyle}>
                    <div style={styles.modalStyle}>
                        <h2>Enter the room name and password</h2>
                        <input
                            type="text"
                            placeholder="Room Name"
                            value={roomName}
                            onChange={(e) => setRoomName(e.target.value)}
                            style={styles.modalInputStyle}
                        />
                        <input
                            type="password"
                            placeholder="Password"
                            value={roomPassword}
                            onChange={(e) => setRoomPassword(e.target.value)}
                            style={styles.modalInputStyle}
                        />
                        <button onClick={handleCreateRoom} style={styles.createButtonStyle}>Create</button>
                        <button onClick={handleModalClose} style={styles.cancelButtonStyle}>Cancel</button>
                    </div>
                </div>
            )}
        </div>
    );
}

function RoomItem({ name, participants }) {
    return (
        <div style={styles.roomItemStyle}>
            <div style={styles.roomAvatarStyle}></div>
            <span style={styles.roomNameStyle}>{name}</span>
            <span>{participants}</span>
        </div>
    );
}

export default HubPage;
