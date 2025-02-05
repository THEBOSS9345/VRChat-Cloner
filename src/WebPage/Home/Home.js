let avatars = [];
let starAvatars = [];

let changingAvatar = false;

let startPage = 0;
const maxAvatarsPerPage = 42;

document.addEventListener('DOMContentLoaded', () => {
    const grid = document.querySelector('.avatar-grid');

    if (!grid) {
        console.error('Grid element not found');
        return;
    }

    const connection = new WebSocket('ws://IP_ADDRESS/avatars');

    connection.onerror = (error) => {
        console.error('Error while connecting to the server:', error);
        showErrorMessage('Failed to connect to the server. Please try again later.');

    };

    connection.onopen = () => {
        console.log('Connected to the server');
    };

    connection.onclose = (event) => {
        console.log('Connection closed');
        if (avatars.length === 0) {
            showErrorMessage('Connection closed. No avatars were loaded.');
        }

    };

    connection.onmessage = (message) => {
        const realm = JSON.parse(message.data);
        avatars.push(realm);

        avatars.sort((a, b) => {
            if (a.star && !b.star) {
                return -1;
            }

            if (!a.star && b.star) {
                return 1;
            }

            if (Date.parse(a.created_at) > Date.parse(b.created_at)) {
                return -1;
            }

            if (Date.parse(a.created_at) < Date.parse(b.created_at)) {
                return 1;
            }

            if (a.equippedTime && b.equippedTime) {
            if (Date.parse(a.equippedTime) > Date.parse(b.equippedTime)) {
                return -1;
            }

            if (Date.parse(a.equippedTime) < Date.parse(b.equippedTime)) {
                return 1;
            }
            }

            return a.name.localeCompare(b.name);
        })

        if (avatars.length <= maxAvatarsPerPage) {
            avatars = document.querySelector('.search').value ? avatars.filter(realm => realm.name.toLowerCase().includes(document.querySelector('.search').value.toLowerCase())) : avatars;

            if (avatars.length <= maxAvatarsPerPage) {
                renderAvatars(avatars);
            } else {
                renderAvatars(avatars.slice(startPage * maxAvatarsPerPage, (startPage + 1) * maxAvatarsPerPage));
            }
        }
    };

    document.addEventListener('click', async (event) => {
        if (event.target.classList.contains('equip-button')) {
            const avatarId = event.target.closest('.avatar-card').querySelector('.avatar-id').textContent;
            // loading spinner

            event.target.closest('.avatar-card').querySelector('.loading-spinner').style.display = 'block';
            showProcessingMessage("Changing avatar...");

            if (changingAvatar) {
                event.target.closest('.avatar-card').querySelector('.loading-spinner').style.display = 'none';
                showErrorMessage('Please wait for your old avatar to finish changing.');
                return
            }

            changingAvatar = true;

            await new Promise(r => setTimeout(r, 4000));


            changingAvatar = true;
            changeAvatar(avatarId, event)
        } else if (event.target.tagName === "SPAN" || event.target.classList.contains('star-icon')) {
            const avatarId = event.target.closest('.avatar-card').querySelector('.avatar-id').textContent;
            await toggleStar(avatarId, event);
        }
    });

  async  function toggleStar(avatarId, event) {
        console.log(`Star button clicked for avatar with ID: ${avatarId}`);

            starAvatars[avatarId] = !starAvatars[avatarId];

            if(changingAvatar) {
                showErrorMessage('Please wait for your old avatar to finish changing.');
                return;
            }


      showProcessingMessage("Changing star status...");
      await new Promise(r => setTimeout(r, 2000));
      event.target.closest('.star-button').classList.add('starred');

            if (starAvatars[avatarId]) {
                StarAvatar(avatarId);
            } else {
                event.target.closest('.star-button').classList.remove('starred');
                StarAvatar(avatarId, false);
            }
    }

    document.getElementById("filter-button").addEventListener('click', () => {
        const search = document.querySelector('.search').value;

        if (search.length === 0) {
            const grid = document.querySelector('.avatar-grid');

            if (!grid) {
                console.error('Grid element not found');
                return;
            }

            if (grid.classList.contains('avatar-card') === false) {
                grid.innerHTML = '';

                startPage = 0;

                updatePagination();
                renderAvatars(avatars.slice(startPage * maxAvatarsPerPage, (startPage + 1) * maxAvatarsPerPage));
                showConfirmMessage("Resetting Avatars...");
                return
            }

            showErrorMessage('Please enter a search query.');
            return;
        }

        grid.innerHTML = '';

        let filteredAvatars = avatars.filter(avatar => avatar.name.toLowerCase().includes(search.toLowerCase()));

        if (filteredAvatars.length === 0) {
            showErrorMessage('No avatars found with that name.');
            return;
        }

        if (filteredAvatars.length <= maxAvatarsPerPage) {
            renderAvatars(filteredAvatars);
        } else {
            startPage = 0;
            updatePagination();
            renderAvatars(filteredAvatars.slice(startPage * maxAvatarsPerPage, (startPage + 1) * maxAvatarsPerPage));
        }
    })

    const previousPageButton = document.getElementById('previous-page');
    const nextPageButton = document.getElementById('next-page');
    const pageButtons = document.querySelectorAll('.pagination-button:not(.nav)');

    previousPageButton.addEventListener('click', () => {
        if (startPage > 0) {
            startPage -= 1;
            updatePagination();
            renderAvatars(avatars.slice(startPage * maxAvatarsPerPage, (startPage + 1) * maxAvatarsPerPage));
        }
    });

    nextPageButton.addEventListener('click', () => {
        if (startPage < Math.floor(avatars.length / maxAvatarsPerPage)) {
            updatePagination();
            startPage += 1;
            renderAvatars(avatars.slice(startPage * maxAvatarsPerPage, (startPage + 1) * maxAvatarsPerPage));
        }
    });





    pageButtons.forEach(button => {
        button.addEventListener('click', () => {

            if (parseInt(button.textContent) > avatars.length / maxAvatarsPerPage) {
                return;
            }

            startPage = parseInt(button.textContent) - 1;

            updatePagination();
            renderAvatars(avatars.slice(startPage * maxAvatarsPerPage, (startPage + 1) * maxAvatarsPerPage));
        });
    });

    function updatePagination() {
        pageButtons.forEach((button, index) => {

            button.textContent = startPage + index + 1;
            button.classList.toggle('active', index === 0);
        });

        previousPageButton.disabled = startPage === 0;
        nextPageButton.disabled = startPage >= Math.floor(avatars.length / maxAvatarsPerPage);
    }

    updatePagination();
});

function StarAvatar(avatarId, star = true) {
    fetch('http://IP_ADDRESS/starAvatar', {
        method: 'POST',
        mode: 'no-cors',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ avatarId, star }),
    })
        .then(async  response => {
            changingAvatar = false;
            showConfirmMessage(star ? "Avatar starred successfully." : "Avatar un-starred successfully.");
        })
        .catch((error) => {
            changingAvatar = false;
            showErrorMessage(`Failed to ${star ? 'star' : 'unstar'} avatar. Please try again.`);
        });
}
function renderAvatars(avatars) {
    const grid = document.querySelector('.avatar-grid');
    if (!grid) {
        console.error('Grid element not found');
        return;
    }

    grid.innerHTML = '';

document.body.style.minHeight = `${Math.ceil(avatars.length / maxAvatarsPerPage) * 275}vh`;

    avatars.forEach(avatar => {
        const card = document.createElement('div');
        card.classList.add('avatar-card');
        card.dataset.id = avatar.id;

        starAvatars[avatar.id] = avatar.star
        card.innerHTML = `
            <div class="avatar-image-wrapper">
                <img src="${avatar.imageUrl}" alt="${avatar.name}" class="avatar-image">
                <button class="star-button">
                    <span class="star-icon">★</span>
                </button>
            </div>
            <div class="avatar-content">
                <h3 class="avatar-name">${avatar.name}</h3>
                <p class="avatar-description">${avatar.description}</p>
                <div>
                <span class="time-label">Saved:</span>
                <span class="time-value">${new Date(avatar.created_at).toLocaleString()}
            </div>
            <div>
                <span class="time-label">Equipped:</span>
                <span class="time-value">${!avatar.equippedTime ? 'Never Used' : new Date(avatar.equippedTime).toLocaleString()}
            </div>
                <div class="avatar-id">${avatar.id}</div>
                <button class="equip-button">Equip</button>
            </div>
            <div id="loading-spinner" class="loading-spinner" style="display: none;"></div>
        `;

        if (starAvatars[avatar.id]) {
            card.querySelector('.star-button').classList.add('starred');
        }

        grid.appendChild(card);
    });
}




function changeAvatar(avatarId, event) {
    fetch('http://IP_ADDRESS/changeAvatar', {
        method: 'POST',
        mode: 'no-cors',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ avatarId }),
    })
        .then(async  response => {
            const FindAvatarIndex = avatars.findIndex(avatar => avatar.id === avatarId);

            avatars[FindAvatarIndex].equippedTime = Date.now();
            event.target.closest('.avatar-card').querySelector('.loading-spinner').style.display = 'none';
            changingAvatar = false;
            return showConfirmMessage("Avatar changed successfully to " + avatars[FindAvatarIndex].name);
        })
        .catch((error) => {
            changingAvatar = false;
            event.target.closest('.avatar-card').querySelector('.loading-spinner').style.display = 'none';
           return showErrorMessage('Failed to change avatar. Please try again.');
        });
}


const messageContainer = document.createElement('div');
messageContainer.classList.add('message-container');
document.body.appendChild(messageContainer);

function showErrorMessage(message) {
    const messageElement = document.createElement('div');
    messageElement.classList.add('message', 'error-message');
    messageElement.textContent = message;
    messageContainer.appendChild(messageElement);
    setTimeout(() => {
        messageElement.remove();
    }, 5000);
}

function showConfirmMessage(message) {
    const messageElement = document.createElement('div');
    messageElement.classList.add('message', 'confirm-message');
    messageElement.textContent = message;
    messageContainer.appendChild(messageElement);
    setTimeout(() => {
        messageElement.remove();
    }, 5000);
}

function showProcessingMessage(message) {
    const messageElement = document.createElement('div');
    messageElement.classList.add('message', 'processing-message');
    messageElement.textContent = message || 'Processing...';
    messageContainer.appendChild(messageElement);
    setTimeout(() => {
        messageElement.remove();
    }, 5000);
}