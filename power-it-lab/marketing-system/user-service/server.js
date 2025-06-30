const express = require('express');
const app = express();
app.use(express.json());

app.get('/health', (req, res) => {
    res.json({ status: 'ok', service: 'user-service' });
});

app.listen(8080, () => {
    console.log('User service running on port 8080');
});
