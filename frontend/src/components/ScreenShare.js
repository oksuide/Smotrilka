import React, { useState } from 'react';
import { socket } from '../socket';

const ScreenShare = () => {
    const [stream, setStream] = useState(null);

    const startScreenShare = async () => {
        try {
            const mediaStream = await navigator.mediaDevices.getDisplayMedia({
                video: true,
                audio: true,
            });
            setStream(mediaStream);
            // Отправка сигнала о начале трансляции
            socket.send(JSON.stringify({ type: 'start_stream', message: 'Screen sharing started' }));
        } catch (err) {
            console.error('Error starting screen share:', err);
        }
    };

    return (
        <div>
            <button onClick={startScreenShare}>Start Screen Share</button>
            {stream && (
                <video
                    id="screen-video"
                    autoPlay
                    playsInline
                    muted
                    style={{ width: '100%' }}
                    srcObject={stream}
                ></video>
            )}
        </div>
    );
};

export default ScreenShare;
