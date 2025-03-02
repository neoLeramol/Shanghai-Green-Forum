// 显示提示信息
function showMessage(message) {
    const messageBox = document.getElementById('message-box');
    messageBox.textContent = message;
    messageBox.style.opacity = 1;
    setTimeout(() => {
        messageBox.style.opacity = 0;
    }, 3000);
}

// function.html 地图初始化
function initMap() {
    var map = new qq.maps.Map(document.getElementById('map'), {
        center: new qq.maps.LatLng(31.230416, 121.473701),
        zoom: 12
    });

    qq.maps.event.addListener(map, 'zoom_changed', function () {
        console.log('地图缩放级别变为: ', map.getZoom());
    });

    qq.maps.event.addListener(map, 'center_changed', function () {
        var center = map.getCenter();
        console.log('地图中心移动到: ', center.lat(), center.lng());
    });

    fetch('/posts')
      .then(response => response.json())
      .then(posts => {
            posts.forEach(post => {
                var position = new qq.maps.LatLng(post.location.lat, post.location.lng);
                var marker = new qq.maps.Marker({
                    position: position,
                    map: map,
                    icon: {
                        url: 'https://map.qq.com/wemap/api/img/marker.png',
                        size: new qq.maps.Size(24, 32),
                        anchor: new qq.maps.Point(12, 32)
                    }
                });

                qq.maps.event.addListener(marker, 'click', function () {
                    showPostInfo(post, map, position);
                });
            });
        });

    return map;
}

// 显示帖子信息
function showPostInfo(post, map, position) {
    var infoWindow = new qq.maps.InfoWindow({
        content: `<div>
                    <h3>帖子编号: ${post.id}</h3>
                    <p>文本内容: ${post.text}</p>
                    <p>日期: ${post.date}</p>
                </div>`
    });
    infoWindow.open(map, position);
}

// master.html 显示搜索或发帖界面
function showSearch() {
    document.getElementById('search-section').style.display = 'block';
    document.getElementById('post-section').style.display = 'none';
}

function showPost() {
    document.getElementById('post-section').style.display = 'block';
    document.getElementById('search-section').style.display = 'none';

    var postMap = new qq.maps.Map(document.getElementById('post-map'), {
        center: new qq.maps.LatLng(31.230416, 121.473701),
        zoom: 12
    });

    var selectedMarker;
    qq.maps.event.addListener(postMap, 'click', function (event) {
        if (selectedMarker) {
            selectedMarker.setMap(null);
        }
        selectedMarker = new qq.maps.Marker({
            position: event.latLng,
            map: postMap
        });
        document.getElementById('post-location').value = `${event.latLng.lat()},${event.latLng.lng()}`;
    });
}

function searchPosts() {
    var id = document.getElementById('search-id').value;
    var date = document.getElementById('search-date').value;
    let url = '/search';
    if (id || date) {
        url += `?id=${id}&date=${date}`;
    }
    fetch(url)
      .then(response => response.json())
      .then(data => {
            var resultsDiv = document.getElementById('search-results');
            resultsDiv.innerHTML = '';
            if (data.length === 0) {
                showMessage('未找到相关帖子');
            }
            data.forEach(post => {
                var postDiv = document.createElement('div');
                postDiv.innerHTML = `编号：${post.id}<br>文本内容：${post.text}<br>日期：${post.date}<br>位置：${post.location}<br>`;
                var deleteButton = document.createElement('button');
                deleteButton.textContent = '删除';
                deleteButton.onclick = function () {
                    deletePost(post.id);
                };
                var modifyButton = document.createElement('button');
                modifyButton.textContent = '修改';
                modifyButton.onclick = function () {
                    modifyPost(post.id);
                };
                postDiv.appendChild(deleteButton);
                postDiv.appendChild(modifyButton);
                resultsDiv.appendChild(postDiv);
            });
        })
      .catch(error => {
            showMessage('搜索出错，请稍后重试');
            console.error(error);
        });
}

function deletePost(id) {
    if (confirm('确定要删除该帖子吗？')) {
        fetch(`/delete?id=${id}`, {
            method: 'DELETE'
        })
          .then(response => response.json())
          .then(data => {
                showMessage(data.message);
                searchPosts();
            })
          .catch(error => {
                showMessage('删除出错，请稍后重试');
                console.error(error);
            });
    }
}

function modifyPost(id) {
    alert(`你要修改编号为 ${id} 的帖子，后续可完善修改逻辑`);
}

function publishPost(location) {
    var text = document.getElementById('post-text').value;
    var image = document.getElementById('post-image').files[0];

    if (!text) {
        showMessage('文本内容不能为空');
        return;
    }

    var formData = new FormData();
    formData.append('text', text);
    formData.append('image', image);
    formData.append('location', location);

    fetch('/publish', {
        method: 'POST',
        body: formData
    })
      .then(response => response.json())
      .then(data => {
            showMessage(data.message);
            if (data.message === '帖子发布成功') {
                document.getElementById('post-text').value = '';
                document.getElementById('post-image').value = '';
                // 重新加载地图标记
                initMap();
            }
        })
      .catch(error => {
            showMessage('发布出错，请稍后重试');
            console.error(error);
        });
}

// 尝试发布帖子，先获取地理位置
function tryPublishPost() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(function (position) {
            const location = `${position.coords.latitude},${position.coords.longitude}`;
            publishPost(location);
        }, function () {
            showMessage('无法获取您的地理位置，请手动选择位置。');
            // 这里可以添加手动选择位置的逻辑
        });
    } else {
        showMessage('您的浏览器不支持地理位置功能，请手动选择位置。');
        // 这里可以添加手动选择位置的逻辑
    }
}

window.onload = function () {
    if (document.getElementById('map')) {
        initMap();
    }
};