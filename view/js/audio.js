var audio = document.getElementById("myAudio");
var playPauseButton = document.getElementById("playPauseButton");
var progressBar = document.getElementById("progressBar");

playPauseButton.addEventListener("click", function() {
  if (audio.paused) {
    audio.play();
    playPauseButton.innerHTML = '<i class="bi bi-pause-fill"></i>'; // Change to play icon
  } else {
    audio.pause();
    playPauseButton.innerHTML = '<i class="bi bi-play-fill"></i>'; // Change to play icon
  }
});

audio.addEventListener("timeupdate", function() {
  var progress = (audio.currentTime / audio.duration) * 100;
  progressBar.style.width = progress + "%";
});