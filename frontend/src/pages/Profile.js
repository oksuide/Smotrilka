import React, { useState } from 'react';

function Profile() {
    const [avatar, setAvatar] = useState(null);
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');

    const handleAvatarChange = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setAvatar(reader.result);
            };
            reader.readAsDataURL(file);
        }
    };

    const handlePasswordChange = () => {
        if (password === confirmPassword) {
            alert('Password changed successfully!');
            setPassword('');
            setConfirmPassword('');
        } else {
            alert('Passwords do not match!');
        }
    };

    return (
        <div style={styles.container}>
            <h2 style={styles.title}>Profile</h2>
            <div style={styles.avatarContainer}>
                {avatar ? (
                    <img src={avatar} alt="Avatar" style={styles.avatar} />
                ) : (
                    <div style={styles.placeholder}>No Avatar</div>
                )}
                <input type="file" accept="image/*" onChange={handleAvatarChange} style={styles.fileInput} />
            </div>
            <div style={styles.passwordContainer}>
                <input
                    type="password"
                    placeholder="New Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    style={styles.input}
                />
                <input
                    type="password"
                    placeholder="Confirm Password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    style={styles.input}
                />
                <button onClick={handlePasswordChange} style={styles.button}>Change Password</button>
            </div>
        </div>
    );
}

const styles = {
    container: { display: 'flex', flexDirection: 'column', alignItems: 'center', padding: '20px', backgroundColor: '#1e1e1e', height: '100vh', color: '#fff' },
    title: { marginBottom: '20px' },
    avatarContainer: { display: 'flex', flexDirection: 'column', alignItems: 'center', marginBottom: '20px' },
    avatar: { width: '100px', height: '100px', borderRadius: '50%', objectFit: 'cover' },
    placeholder: { width: '100px', height: '100px', borderRadius: '50%', backgroundColor: '#444', display: 'flex', alignItems: 'center', justifyContent: 'center' },
    fileInput: { marginTop: '10px' },
    passwordContainer: { display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '10px' },
    input: { padding: '10px', borderRadius: '5px', border: 'none', width: '200px' },
    button: { padding: '10px 20px', borderRadius: '5px', border: 'none', backgroundColor: '#444', color: '#fff', cursor: 'pointer' }
};

export default Profile;
