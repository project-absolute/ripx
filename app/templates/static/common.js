document.addEventListener('DOMContentLoaded', function () {
  const uploadArea = document.getElementById('uploadArea');
  const fileInput = document.getElementById('fileInput');
  const uploadForm = document.getElementById('uploadForm') || document.getElementById('imageUploadForm');

  if (uploadArea && fileInput && uploadForm) {
    uploadArea.addEventListener('click', function () { fileInput.click() });
    fileInput.addEventListener('change', function () {
      if (fileInput.files.length > 0) {
        uploadForm.submit();
      }
    });
    uploadArea.addEventListener('dragover', function (e) { e.preventDefault(); uploadArea.classList.add('dragover') });
    uploadArea.addEventListener('dragleave', function (e) { e.preventDefault(); uploadArea.classList.remove('dragover') });
    uploadArea.addEventListener('drop', function (e) {
      e.preventDefault(); uploadArea.classList.remove('dragover');
      const files = e.dataTransfer.files;
      if (files.length > 0) {
        fileInput.files = files;
        uploadForm.submit();
      }
    });
  }
});

function copyUrl(sessionID, albumID, filename, button) {
  const url = window.location.origin + '/' + sessionID + '/' + albumID + '/' + filename;
  if (navigator.clipboard) {
    navigator.clipboard.writeText(url).then(function () {
      const originalText = button.textContent;
      button.textContent = 'ᴄᴋоᴨиᴩоʙᴀно!';
      button.classList.add('copied');
      setTimeout(function () { button.textContent = originalText; button.classList.remove('copied') }, 2000);
    }).catch(function (err) { console.error('нᴇ удᴀᴧоᴄь ᴄᴋоᴨиᴩоʙᴀᴛь ᴜʀʟ: ', err) });
  } else {
    const textArea = document.createElement('textarea');
    textArea.value = url;
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    try {
      document.execCommand('copy');
      const originalText = button.textContent;
      button.textContent = 'ᴄᴋоᴨиᴩоʙᴀно!';
      button.classList.add('copied');
      setTimeout(function () { button.textContent = originalText; button.classList.remove('copied') }, 2000);
    } catch (err) { console.error('Не удалось скопировать URL: ', err) }
    document.body.removeChild(textArea);
  }
}

function deleteImage(sessionID, albumID, filename, button) {
  if (!confirm('Вы уверены, что хотите удалить это изображение?')) {
    return;
  }

  const formData = new FormData();
  formData.append('album_id', albumID);
  formData.append('filename', filename);

  fetch('/delete-image', {
    method: 'POST',
    body: formData
  })
    .then(response => {
      if (response.ok) {
        // Удаляем элемент изображения из DOM
        const imageItem = button.closest('.image-item');
        imageItem.style.opacity = '0.5';
        setTimeout(() => {
          imageItem.remove();
          // Проверяем, остались ли изображения
          const remainingImages = document.querySelectorAll('.image-item');
          if (remainingImages.length === 0) {
            // Показываем пустое состояние
            const imageGrid = document.getElementById('imageGrid');
            imageGrid.innerHTML = '<div class="empty-state"><div class="empty-icon">📷</div><div class="empty-text">у ʙᴀᴄ ᴨоᴋᴀ нᴇᴛ зᴀᴦᴩужᴇнных изобᴩᴀжᴇний</div><a href="/" class="empty-link">зᴀᴦᴩузиᴛь ᴨᴇᴩʙоᴇ изобᴩᴀжᴇниᴇ</a></div>';
          }
        }, 300);
      } else {
        alert('Ошибка при удалении изображения');
      }
    })
    .catch(error => {
      console.error('Error:', error);
      alert('Ошибка при удалении изображения');
    });
}

function deleteUser() {
  if (!confirm('Вы уверены, что хотите удалить весь профиль со всеми альбомами и изображениями? Это действие необратимо!')) {
    return;
  }

  fetch('/delete-user', {
    method: 'POST'
  })
    .then(response => {
      if (response.ok) {
        // Удаляем cookie на клиенте
        document.cookie = 'session_id=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
        // Перезагружаем страницу для получения новой сессии
        window.location.reload();
      } else {
        alert('Ошибка при удалении профиля');
      }
    })
    .catch(error => {
      console.error('Error:', error);
      alert('Ошибка при удалении профиля');
    });
}
