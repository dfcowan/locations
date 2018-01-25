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
        center: home,
        styles: []
    });
}

function heatUp() {
    var userId = getQuerystringParameterValue("user");

    $.ajax({
        url: '/api/users/' + userId + '/counts'
    }).done(function (data, textStatus, jqXHR) {
        var heatMapData = [];
        var usersCenter = {
            location: home,
            weight: 0
        };

        data.forEach(function(cc) {
            var l = {
                location: new google.maps.LatLng(
                    cc.Coordinate.Lat,
                    cc.Coordinate.Lon),
                weight: cc.Count	
            };
            heatMapData.push(l);
            if(l.weight > usersCenter.weight) {
                usersCenter = l;
            }
        });

        map.setCenter(usersCenter.location);
            
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
    }).fail(function (jqXHR, textStatus, errorThrown) {
        alert('failed to load counts');
    });
        
    $.ajax({
        url: '/api/users/' + userId + '/sync'
    }).done(function (data, textStatus, jqXHR) {
        var startDate = addHyphens(data.StartDate);
        var endDate = addHyphens(data.SyncedThroughDate);
        
        $("#startDate").val(startDate);
        $("#startDate").attr("min", startDate);
        $("#endDate").attr("min", startDate);

        applyEndDate(endDate);
    }).fail(function(){
        alert('failed to load sync status');
    });    
}

function addHyphens(d) {
    return d.substr(0,4) + "-" + d.substr(4,2) + "-" + d.substr(6,2);
}

function removeHyphen(d) {
    return d.substr(0,4) + d.substr(5,2) + d.substr(8,2);
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
    var startDate = removeHyphen($("#startDate").val());
    var endDate = removeHyphen($("#endDate").val());

    $.ajax({
        url: '/api/users/' + userId + '/counts?startDate=' + startDate + '&endDate=' + endDate
    }).done(function (data, textStatus, jqXHR) {
        var heatMapData = [];
        data.forEach(function(cc) {
            heatMapData.push(
                {
                    location: new google.maps.LatLng(
                        cc.Coordinate.Lat,
                        cc.Coordinate.Lon),
                        weight: cc.Count	
                }
            );
        });

        heatmap.setData(heatMapData);

        $("#refresh").text("Limit by Date Range");
    }).fail(function (jqXHR, textStatus, errorThrown) {
        alert('failed to load counts');
    });    
}

function sync() {
    var userId = getQuerystringParameterValue("user");

    $("#sync").text("Syncing...");

    $.ajax({
        url: '/api/users/' + userId + '/sync',
        method: "POST"
    }).done(function (data, textStatus, jqXHR) {
        var endDate = addHyphens(data.SyncedThroughDate);
        endDate = applyEndDate(endDate);
        $("#sync").text("Synced through " + endDate);
    }).fail(function(){
        alert('failed to sync');
    });    
}

function applyEndDate(endDate) {
    var dt = new Date(endDate);
    dt.setDate(dt.getDate() - 1);
    dt = dt.toISOString().substr(0,10);

    if (!$("#endDate").val() || $("#endDate").val() == $("#endDate").attr("max")) {
        $("#endDate").val(dt);
    }

    $("#startDate").attr("max", dt);
    $("#endDate").attr("max", dt);

    return dt;
}

$(document).ready(function () {
    heatUp();
});