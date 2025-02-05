import React, { useRef, useState } from 'react';
import { FaMicrophone, FaMicrophoneSlash, FaTv, FaVolumeMute, FaVolumeUp } from "react-icons/fa";
import Chat from '/home/oleg/goprojects/Smotrilka/frontend/src/components/Chat';

function Room() {
    const [isSharing, setIsSharing] = useState(false);
    const [isMuted, setIsMuted] = useState(false);
    const [isMicOff, setIsMicOff] = useState(false);
    const [activeTab, setActiveTab] = useState('chat');
    const [isChatVisible, setIsChatVisible] = useState(true);
    const videoRef = useRef(null);

    const startScreenShare = async () => {
        try {
            const stream = await navigator.mediaDevices.getDisplayMedia({ video: true });
            if (videoRef.current) {
                videoRef.current.srcObject = stream;
            }
            setIsSharing(true);
        } catch (error) {
            console.error('Ошибка при выборе экрана:', error);
        }
    };


    const switchScreen = async () => {
        await startScreenShare();
    };

    return (
        <div style={styles.container}>
            <div style={{ ...styles.videoContainer, flex: isChatVisible ? 2 : 3 }}>
                {isSharing ? (
                    <video ref={videoRef} autoPlay style={styles.video} />
                ) : (
                    <button style={styles.shareButton} onClick={startScreenShare}>
                        Start share your screen
                    </button>
                )}
                {isChatVisible && (
                    <div style={styles.chatSidebar}>
                        <button style={styles.hideChatButton} onClick={() => setIsChatVisible(false)}>▶</button>
                    </div>
                )}
                <div style={styles.chatSidebar}>
                    {isChatVisible && (
                        <button style={styles.hideChatButton} onClick={() => setIsChatVisible(false)}>▶</button>
                    )}
                </div>
            </div>
            <div style={styles.controls}>
                <button onClick={() => setIsMicOff(!isMicOff)} style={styles.controlButton}>
                    {isMicOff ? <FaMicrophoneSlash /> : <FaMicrophone />}
                </button>
                <button onClick={() => setIsMuted(!isMuted)} style={styles.controlButton}>
                    {isMuted ? <FaVolumeMute /> : <FaVolumeUp />}
                </button>
                <button onClick={switchScreen} style={styles.controlButton}>
                    {<FaTv />}</button>
            </div>

            {isChatVisible && <Chat activeTab={activeTab} setActiveTab={setActiveTab} setIsChatVisible={setIsChatVisible} />}

            {!isChatVisible && (
                <button style={styles.showChatButton} onClick={() => setIsChatVisible(true)}>◀</button>
            )}
        </div>
    );
}

const styles = {
    container: { display: 'flex', height: '100vh', backgroundColor: '#1e1e1e', position: 'relative' },
    videoContainer: { display: 'flex', justifyContent: 'center', alignItems: 'center', color: '#fff', transition: 'flex 0.3s ease' },
    video: { width: '100%', height: '100%', objectFit: 'contain' },
    shareButton: { padding: '15px', backgroundColor: '#444', color: '#fff', border: 'none', cursor: 'pointer' },
    hideChatButton: { position: 'absolute', top: '10px', right: '300px', background: 'transparent', color: '#666', border: 'none', cursor: 'pointer' },
    showChatButton: { position: 'absolute', top: '10px', right: '10px', background: 'transparent', color: '#666', border: 'none', cursor: 'pointer' },
    controls: { position: 'absolute', bottom: '10px', left: '10px', display: 'flex', gap: '10px' },
    controlButton: { width: '40px', height: '40px', borderRadius: '50%', backgroundColor: '#fff', color: '#000', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', justifyContent: 'center' }
};

export default Room;
