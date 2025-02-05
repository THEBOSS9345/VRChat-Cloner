const canvas = document.getElementById('constellation-bg');
const ctx = canvas.getContext('2d');

let width, height, points;
let mouse = { x: 0, y: 0 };

function init() {
    width = canvas.width = window.innerWidth;
    height = canvas.height = window.innerHeight;

    points = [];
    for (let i = 0; i < 100; i++) {
        points.push({
            x: Math.random() * width,
            y: Math.random() * height,
            vx: (Math.random() - 0.5) * 0.5,
            vy: (Math.random() - 0.5) * 0.5
        });
    }
}

function animate() {
    ctx.clearRect(0, 0, width, height);

    points.forEach(point => {
        point.x += point.vx;
        point.y += point.vy;


        if (point.x < 0 || point.x > width) point.vx *= -1;
        if (point.y < 0 || point.y > height) point.vy *= -1;


        let dx = mouse.x - point.x;
        let dy = mouse.y - point.y;
        let distance = Math.sqrt(dx * dx + dy * dy);
        if (distance < 100) {
            let angle = Math.atan2(dy, dx);
            point.vx -= Math.cos(angle) * 0.02;
            point.vy -= Math.sin(angle) * 0.02;
        }
        ctx.beginPath();
        ctx.arc(point.x, point.y, 1, 0, Math.PI * 2);
        ctx.fillStyle = 'rgba(0, 128, 255, 0.5)';
        ctx.fill();

        points.forEach(otherPoint => {
            const distance = Math.sqrt(
                Math.pow(point.x - otherPoint.x, 2) + Math.pow(point.y - otherPoint.y, 2)
            );
            if (distance < 100) {
                ctx.beginPath();
                ctx.moveTo(point.x, point.y);
                ctx.lineTo(otherPoint.x, otherPoint.y);
                ctx.strokeStyle = `rgba(255, 255, 255, ${0.2 - distance / 500})`;
                ctx.stroke();
            }
        });
    });

    requestAnimationFrame(animate);
}

window.addEventListener('resize', init);
window.addEventListener('mousemove', (e) => {
    mouse.x = e.clientX;
    mouse.y = e.clientY;
});

init();
animate();

document.getElementById('login-form').addEventListener('submit', function(e) {
    e.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    console.log('Login attempted with:', { email, password });
    // Here you would typically send the login data to your server
});