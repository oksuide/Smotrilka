import React, { useEffect, useState } from 'react';

const socket = new WebSocket('ws://localhost:5000');

function Chat() {
    const [activeTab, setActiveTab] = useState('chat');
    const [messages, setMessages] = useState(() => {
        return JSON.parse(localStorage.getItem('chatMessages')) || [];
    });
    const [inputValue, setInputValue] = useState('');

    useEffect(() => {
        const socket = new WebSocket('ws://localhost:5000');

        socket.onmessage = (event) => {
            try {
                const parsedData = JSON.parse(event.data);
                setMessages((prevMessages) => {
                    const newMessages = [...prevMessages, parsedData.text];
                    localStorage.setItem('chatMessages', JSON.stringify(newMessages));
                    return newMessages;
                });
            } catch (error) {
                console.error('Ошибка обработки входящего сообщения:', error);
            }
        };

        return () => {
            socket.close();
        };
    }, []);

    const sendMessage = () => {
        if (inputValue.trim()) {
            const messageObject = { text: inputValue };
            socket.send(JSON.stringify(messageObject));
            setInputValue('');
        }
    };

    return (
        <div style={styles.chatSidebar}>
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
                        <div style={styles.chatMessages}>
                            {messages.map((msg, index) => (
                                <div key={index} style={styles.message}>{msg}</div>
                            ))}
                        </div>
                        <div style={styles.inputContainer}>
                            <textarea
                                placeholder="Type a message..."
                                style={styles.inputBox}
                                rows="1"
                                value={inputValue}
                                onChange={(e) => setInputValue(e.target.value)}
                                onInput={(e) => {
                                    e.target.style.height = 'auto';
                                    e.target.style.height = e.target.scrollHeight + 'px';
                                }}
                            ></textarea>
                            <button style={styles.sendButton} onClick={sendMessage}>Send</button>
                        </div>
                    </div>
                ) : (
                    <p>Members list...</p>
                )}
            </div>
        </div>
    );
}

const styles = {
    chatSidebar: { width: '300px', backgroundColor: '#222', color: '#fff', display: 'flex', flexDirection: 'column', position: 'relative' },
    tabSwitcher: { display: 'flex' },
    tabButton: { flex: 1, padding: '10px', border: 'none', cursor: 'pointer', color: '#fff' },
    tabContent: { flex: 1, display: 'flex', flexDirection: 'column', justifyContent: 'space-between', padding: '15px' },
    chatBox: { display: 'flex', flexDirection: 'column', flex: 1, maxWidth: '250px' },
    chatMessages: { flex: 1, overflowY: 'auto', maxHeight: '400px' }, // Контейнер сообщений занимает всё доступное пространство
    inputContainer: { display: 'flex', flexDirection: 'column', gap: '5px', position: 'sticky', bottom: '0', backgroundColor: '#222', padding: '10px' }, // Бокс ввода и кнопка остаются внизу
    inputBox: { width: '100%', padding: '6px', border: 'none', resize: 'none', minHeight: '20px', maxHeight: '100px', overflowY: 'hidden', backgroundColor: '#444', color: '#fff' },
    sendButton: { padding: '10px', backgroundColor: '#444', color: '#fff', border: 'none', cursor: 'pointer', maxWidth: '50px', alignSelf: 'flex-end' },
    message: { padding: '5px', backgroundColor: '#333', margin: '5px 0', borderRadius: '5px' }
};
export default Chat;
