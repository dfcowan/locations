function getQuerystringParameterValue(name, url) {
    if (!url) {
        url = window.location.href;
    }
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
    results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

var map;
var heatmap;
var home = {lat: 29.576679, lng: -98.450644};

function initMap() {
    map = new google.maps.Map(document.getElementById('map'), {
        zoom: 12,
        gestureHandling: 'greedy',
        center: home,
        styles: []
    });

    var heatMapData = [];
    heatmap = new google.maps.visualization.HeatmapLayer({
        data: heatMapData,
        gradient: [
            'rgba(0, 255, 255, 0)',
            'rgba(0, 255, 255, 1)',
            'rgba(0, 191, 255, 1)',
            'rgba(0, 127, 255, 1)',
            'rgba(0, 63, 255, 1)',
            'rgba(0, 0, 255, 1)',
            'rgba(0, 0, 223, 1)',
            'rgba(0, 0, 191, 1)',
            'rgba(0, 0, 159, 1)',
            'rgba(0, 0, 127, 1)',
            'rgba(63, 0, 91, 1)',
            'rgba(127, 0, 63, 1)',
            'rgba(191, 0, 31, 1)',
            'rgba(255, 0, 0, 1)'
        ],
        maxIntensity: 15,
        opacity: 1,
        radius: 5
    });
    heatmap.setMap(map);

    heatUp();
}

function heatUp() {
    var userId = getQuerystringParameterValue("user");

    $.ajax({
        url: '/api/users/' + userId + '/sync'
    }).done(function (data, textStatus, jqXHR) {
        var startDate = data.startDate;
        var endDate = data.syncedThroughDate;

        $("#startDate").attr("min", startDate);
        $("#endDate").attr("min", startDate);

        endDate = new Date(endDate);
        endDate.setDate(endDate.getDate());
        endDate = endDate.toISOString().substr(0,10);

        if (!$("#endDate").val() || $("#endDate").val() == $("#endDate").attr("max")) {
            $("#endDate").val(endDate);
        }

        $("#startDate").attr("max", endDate);
        $("#endDate").attr("max", endDate);

        startDate = new Date(endDate)
        startDate.setDate(startDate.getDate() - 6);
        startDate = startDate.toISOString().substr(0,10);

        $("#startDate").val(startDate);

        refresh();
    }).fail(function(){
        alert('failed to load sync status');
    });
}

function changeDateRangeToday() {
    let today = new Date().toISOString().substr(0,10);

    $("#endDate").attr("max", today);
    $("#endDate").val(today);

    $("#startDate").val(today);

    refresh();
}

function changeDateRangeThisWeek() {
    let today = new Date();

    $("#endDate").attr("max", today.toISOString().substr(0,10));
    $("#endDate").val(today.toISOString().substr(0,10));

    let startDate = today
    startDate.setDate(startDate.getDate() - startDate.getDay());

    $("#startDate").val(startDate.toISOString().substr(0, 10));

    refresh();
}

function changeDateRangeThisMonth() {
    let today = new Date();

    $("#endDate").attr("max", today.toISOString().substr(0,10));
    $("#endDate").val(today.toISOString().substr(0,10));

    let startDate = today
    startDate.setDate(1);

    $("#startDate").val(startDate.toISOString().substr(0, 10));

    refresh();
}

function changeDateRangeThisYear() {
    let today = new Date();

    $("#endDate").attr("max", today.toISOString().substr(0,10));
    $("#endDate").val(today.toISOString().substr(0,10));

    let startDate = today
    startDate.setMonth(0, 1);

    $("#startDate").val(startDate.toISOString().substr(0, 10));

    refresh();
}

function changeDateRangeLast7() {
    let today = new Date();

    $("#endDate").attr("max", today.toISOString().substr(0,10));
    $("#endDate").val(today.toISOString().substr(0,10));

    let startDate = today
    startDate.setDate(startDate.getDate() - 6);

    $("#startDate").val(startDate.toISOString().substr(0, 10));

    refresh();
}

function changeDateRangeLast30() {
    let today = new Date();

    $("#endDate").attr("max", today.toISOString().substr(0,10));
    $("#endDate").val(today.toISOString().substr(0,10));

    let startDate = today
    startDate.setDate(startDate.getDate() - 29);

    $("#startDate").val(startDate.toISOString().substr(0, 10));

    refresh();
}

function changeDateRangeLast365() {
    let today = new Date();

    $("#endDate").attr("max", today.toISOString().substr(0,10));
    $("#endDate").val(today.toISOString().substr(0,10));

    let startDate = today
    startDate.setDate(startDate.getDate() - 364);

    $("#startDate").val(startDate.toISOString().substr(0, 10));

    refresh();
}

function changeRadius() {
    var radius = $("#radius").val();
    heatmap.set('radius', radius);
}

function changeMaxIntensity() {
    var maxIntensity = $("#maxIntensity").val();
    heatmap.set('maxIntensity', maxIntensity);
}

function refresh() {
    $("#refresh").text("Loading...");

    var userId = getQuerystringParameterValue("user");
    var startDate = $("#startDate").val();
    var endDate = $("#endDate").val();

    $.ajax({
        url: '/api/users/' + userId + '/counts?startDate=' + startDate + '&endDate=' + endDate
    }).done(function (data, textStatus, jqXHR) {
        var heatMapData = [];

        data.forEach(function(cc) {
            var l = {
                location: new google.maps.LatLng(
                    cc.latitude,
                    cc.longitude),
                weight: cc.count
            };
            heatMapData.push(l);
        });

        heatmap.setData(heatMapData);

        $("#refresh").text("Limit by Date Range");
    }).fail(function (jqXHR, textStatus, errorThrown) {
        alert('failed to load counts');
    });
}