import React from 'react';
import { Link } from 'react-router-dom';

function Login() {
    return (
        <div style={backgroundStyle}>
            <div style={containerStyle}>
                <div style={titleContainerStyle}>
                    <h1 style={titleStyle}>Название</h1>
                </div>
                <form style={formStyle} onSubmit={handleSignIn}>
                    <div style={formGroupStyle}>
                        <label htmlFor="username" style={labelStyle}>Username</label>
                        <input type="text" id="username" name="username" placeholder="Enter username" style={inputStyle} />
                    </div>
                    <div style={formGroupStyle}>
                        <label htmlFor="password" style={labelStyle}>Password</label>
                        <input type="password" id="password" name="password" placeholder="Enter password" style={inputStyle} />
                    </div>
                    <div style={buttonContainerStyle}>
                        <button type="submit" style={signInButtonStyle}>Sign In</button>
                        <Link to="/register">
                            <button type="button" style={registerButtonStyle}>Register</button>
                        </Link>
                    </div>
                </form>
            </div>
        </div>
    );
}

const handleSignIn = (event) => {
    event.preventDefault();
    const username = event.target.username.value;
    const password = event.target.password.value;
    console.log('Sign In clicked');
    console.log('Username:', username);
    console.log('Password:', password);
};

// Новый стиль для фона страницы
const backgroundStyle = {
    backgroundColor: '#1e1e1e',
    height: '100vh',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    margin: 0,
};

// Обновлённый стиль для контейнера
const containerStyle = {
    textAlign: 'center',
    padding: '30px',
    border: '2px solid #ffffff',
    borderRadius: '15px',
    width: '400px',
    backgroundColor: '#2a2a2a',
    color: '#ffffff',
    boxShadow: '0 0 15px rgba(0, 0, 0, 0.7)',
};

// Новый стиль для контейнера заголовка
const titleContainerStyle = {
    marginBottom: '25px',
};

const titleStyle = {
    fontSize: '28px',
    color: '#FFFFFF',
    margin: '0',
};

const formStyle = {
    textAlign: 'left',
};

const formGroupStyle = {
    marginBottom: '20px',
};

const labelStyle = {
    display: 'block',
    marginBottom: '5px',
    fontWeight: 'bold',
};

const inputStyle = {
    width: '100%',
    padding: '10px',
    borderRadius: '5px',
    border: '1px solid #cccccc',
    boxSizing: 'border-box',
    marginBottom: '10px',
};

const buttonContainerStyle = {
    flexDirection: 'column',
    alignItems: 'center',
    marginTop: '20px',
};

const signInButtonStyle = {
    width: '100%',
    padding: '12px 25px',
    border: 'none',
    borderRadius: '5px',
    backgroundColor: '#28a745', // Зеленый для Sign In
    color: '#ffffff',
    cursor: 'pointer',
    margin: '10px 0',
    fontWeight: 'bold',
};

const registerButtonStyle = {
    width: '100%',
    padding: '12px 25px',
    border: 'none',
    borderRadius: '5px',
    backgroundColor: '#007bff', // Синий для Register
    color: '#ffffff',
    cursor: 'pointer',
    margin: '10px 0',
    fontWeight: 'bold',
};

export default Login;
