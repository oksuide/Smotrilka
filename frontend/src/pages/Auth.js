import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from '../styles/AuthPageStyles';

function Auth() {
    const [view, setView] = useState('main');
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleLogin = async () => {
        setError('');
        const response = await fetch('http://localhost:8080/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();
        if (response.ok) {
            localStorage.setItem('token', data.token);
            navigate('/hub');
        } else {
            setError(data.error || 'Ошибка входа');
        }
    };

    const handleRegister = async () => {
        if (password !== confirmPassword) {
            setError('Пароли не совпадают');
            return;
        }

        setError('');
        const response = await fetch('http://localhost:8080/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();
        if (response.ok) {
            setView('login');
        } else {
            setError(data.error || 'Ошибка регистрации');
        }
    };

    return (
        <div style={styles.background}>
            <div style={styles.container}>
                {view === 'main' && (
                    <>
                        <h2 style={styles.title}>Welcome to Smotrilka</h2>
                        <button style={styles.button} onClick={() => setView('login')}>Login</button>
                        <p style={styles.text}>First time?</p>
                        <button style={styles.button} onClick={() => setView('register')}>Register</button>
                    </>
                )}

                {(view === 'login' || view === 'register') && (
                    <div style={styles.fadeIn}>
                        <h2 style={styles.title}>{view === 'login' ? 'Sign In' : 'Register'}</h2>
                        {error && <p style={styles.error}>{error}</p>}
                        <input type="text" placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} style={styles.input} />
                        <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} style={styles.input} />
                        {view === 'register' && (
                            <input type="password" placeholder="Confirm Password" value={confirmPassword} onChange={(e) => setConfirmPassword(e.target.value)} style={styles.input} />
                        )}
                        <button style={styles.button} onClick={view === 'login' ? handleLogin : handleRegister}>
                            {view === 'login' ? 'Sign In' : 'Register'}
                        </button>
                        <p style={styles.link} onClick={() => setView(view === 'login' ? 'register' : 'login')}>
                            {view === 'login' ? 'Back to Register' : 'Back to Login'}
                        </p>
                    </div>
                )}
            </div>
        </div>
    );
}

export default Auth;