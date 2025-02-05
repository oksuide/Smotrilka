const styles = {
    background: {
        backgroundColor: '#1e1e1e',
        height: '100vh',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        margin: 0,
    },
    container: {
        textAlign: 'center',
        padding: '30px',
        border: '2px solid #ffffff',
        borderRadius: '15px',
        width: '400px',
        backgroundColor: '#2a2a2a',
        color: '#ffffff',
        boxShadow: '0 0 15px rgba(0, 0, 0, 0.7)',
    },
    title: {
        fontSize: '24px',
        marginBottom: '20px',
    },
    text: {
        margin: '10px 0',
    },
    button: {
        width: '100%',
        padding: '12px 25px',
        border: 'none',
        borderRadius: '5px',
        backgroundColor: '#007bff',
        color: '#ffffff',
        cursor: 'pointer',
        margin: '10px 0',
        fontWeight: 'bold',
    },
    input: {
        width: '100%',
        padding: '10px',
        margin: '10px 0',
        borderRadius: '5px',
        border: '1px solid #cccccc',
        boxSizing: 'border-box',
    },
    link: {
        color: '#007bff',
        cursor: 'pointer',
        textDecoration: 'underline',
        marginTop: '10px',
    },
    error: {
        color: 'red',
        marginBottom: '10px',
    },
    fadeIn: {
        animation: 'fadeIn 0.5s',
    },
};
export default styles;