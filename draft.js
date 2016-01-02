//set by template
//var numPacks = 3;
//var numPicks = [1, 15, 15, 15];

var currentPack = 0;
var currentPick = 1;
var navHistory = [];
for (var i =0; i<=numPacks; i++) {
	navHistory.push(1);
}

function displayPick(pack, pick) {
	//alert("display" + pack + " " + pick);
	saveNavHistory();
	$(currentPageId()).hide();	
	currentPack = pack;
	currentPick = pick;
	$(currentPageId()).show();
}

function saveNavHistory() {
	if(currentPick == numPicks[currentPack]) {
		navHistory[currentPack] = 1;
	} else {
		navHistory[currentPack] = currentPick;
	}
}

function nextPick() {		
	var nextPack = currentPack;
	var nextPick = currentPick + 1;
	if (nextPick > numPicks[currentPack]) {
		nextPick = 1;
		nextPack += 1;
		if (nextPack > numPacks) {
			nextPack = 0;
		}
	}
	displayPick(nextPack, nextPick);
}

function prevPick() {
	var nextPack = currentPack;
	var nextPick = currentPick - 1;
	if (nextPick < 1) {
		nextPack -= 1;
		if (nextPack < 0) {
			nextPack = numPacks;
		}
		nextPick = numPicks[nextPack];
	}
	displayPick(nextPack, nextPick);
}

function nextPage() {
	var nextPack = currentPack + 1;
	if(nextPack > numPacks) {
		//nextPack = 0;
		return;
	}
	displayPick(nextPack, navHistory[nextPack]);
}

function prevPage() {
	var nextPack = currentPack - 1;
	if (nextPack < 0) {
		//nextPack = numPacks;
		return;
	}
	displayPick(nextPack, navHistory[nextPack]);
}

function currentPageId() {
	if (currentPack == 0) {
		return '#deck';
	} else {
		return '#P' + currentPack + 'p' + currentPick;
	}
}

function atTop() {
	return $(window).scrollTop() == 0;
}

function atBottom() {
	return $(window).scrollTop() + $(window).height() == $(document).height();
}

$(function () {
	$('#rightnav').click(nextPick);
	$('#leftnav').click(prevPick);
	$(document).keydown(function (e) {
		if(e.which == 39) {
			nextPick();
		} else if(e.which == 37) {
			prevPick();
		} else if(e.which == 40 && atBottom()) {
			nextPage();
		} else if(e.which == 38 && atTop()) {
			prevPage();
		}
	});
	$('.packlink').click(function (e) {
		var pack = $(e.target).data('pack') + 1;
		displayPick(pack, 1);
		return false;
	});
});

