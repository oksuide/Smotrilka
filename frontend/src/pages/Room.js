import React, { useState } from 'react';

function Room() {
    const [isSharing, setIsSharing] = useState(false);
    const [activeTab, setActiveTab] = useState('chat'); // 'chat' или 'members'
    const [isChatVisible, setIsChatVisible] = useState(true);

    const startScreenShare = async () => {
        try {
            await navigator.mediaDevices.getDisplayMedia({ video: true });
            setIsSharing(true);
        } catch (error) {
            console.error('Ошибка при выборе экрана:', error);
        }
    };

    return (
        <div style={styles.container}>
            <div style={styles.videoContainer}>
                {isSharing ? (
                    <p>Screen is being shared...</p>
                ) : (
                    <button style={styles.shareButton} onClick={startScreenShare}>
                        Start share your screen
                    </button>
                )}
            </div>

            <div style={styles.controls}>
                <button style={styles.controlButton}>🎤</button>
                <button style={styles.controlButton}>🔊</button>
            </div>

            {isChatVisible && (
                <div style={styles.chatSidebar}>
                    <button style={styles.hideChatButton} onClick={() => setIsChatVisible(false)}>▶</button>
                    <div style={styles.tabSwitcher}>
                        <button
                            style={{ ...styles.tabButton, backgroundColor: activeTab === 'chat' ? '#444' : '#222' }}
                            onClick={() => setActiveTab('chat')}
                        >
                            Chat
                        </button>
                        <button
                            style={{ ...styles.tabButton, backgroundColor: activeTab === 'members' ? '#444' : '#222' }}
                            onClick={() => setActiveTab('members')}
                        >
                            Members
                        </button>
                    </div>
                    <div style={styles.tabContent}>
                        {activeTab === 'chat' ? (
                            <div style={styles.chatBox}>
                                <div style={styles.chatMessages}>Chat messages...</div>
                                <div style={styles.inputContainer}>
                                    <textarea
                                        placeholder="Type a message..."
                                        style={styles.inputBox}
                                        rows="1"
                                        onInput={(e) => {
                                            e.target.style.height = 'auto';
                                            e.target.style.height = e.target.scrollHeight + 'px';
                                        }}
                                    ></textarea>
                                    <button style={styles.sendButton}>Send</button>
                                </div>
                            </div>
                        ) : (
                            <p>Members list...</p>
                        )}
                    </div>
                </div>
            )}

            {!isChatVisible && (
                <button style={styles.showChatButton} onClick={() => setIsChatVisible(true)}>◀</button>
            )}
        </div>
    );
}

const styles = {
    container: { display: 'flex', height: '100vh', backgroundColor: '#1e1e1e', position: 'relative' },
    videoContainer: { flex: 1, display: 'flex', justifyContent: 'center', alignItems: 'center', color: '#fff' },
    shareButton: { padding: '15px', backgroundColor: '#444', color: '#fff', border: 'none', cursor: 'pointer' },
    chatSidebar: { width: '300px', backgroundColor: '#222', color: '#fff', display: 'flex', flexDirection: 'column', position: 'relative' },
    tabSwitcher: { display: 'flex' },
    tabButton: { flex: 1, padding: '10px', border: 'none', cursor: 'pointer', color: '#fff' },
    tabContent: { flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'space-between', padding: '15px' },
    chatBox: { display: 'flex', flexDirection: 'column', flex: 1 },
    chatMessages: { flex: 1, overflowY: 'auto' },
    inputContainer: { display: 'flex', flexDirection: 'column', gap: '5px' },
    inputBox: { width: '100%', padding: '10px', border: 'none', resize: 'none', minHeight: '40px', maxHeight: '200px', overflowY: 'hidden' },
    sendButton: { padding: '10px', backgroundColor: '#444', color: '#fff', border: 'none', cursor: 'pointer' },
    hideChatButton: { position: 'absolute', top: '10px', left: '-20px', background: 'transparent', color: '#666', border: 'none', cursor: 'pointer' },
    showChatButton: { position: 'absolute', top: '10px', right: '10px', background: 'transparent', color: '#666', border: 'none', cursor: 'pointer' },
    controls: { position: 'absolute', bottom: '10px', left: '10px', display: 'flex', gap: '10px' },
    controlButton: { width: '40px', height: '40px', borderRadius: '50%', backgroundColor: '#fff', color: '#000', border: 'none', cursor: 'pointer', display: 'flex', alignItems: 'center', justifyContent: 'center' }
};

export default Room;