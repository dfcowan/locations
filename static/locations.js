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

function initMap() {
    console.log('initMap');
    
    var home = {lat: 29.576679, lng: -98.450644};
    
    map = new google.maps.Map(document.getElementById('map'), {
        zoom: 12,
        center: home,
        styles: []
    });
    
    console.log('now you see me');
}

function heatUp() {
    console.log('heatUp');
    
    var userId = getQuerystringParameterValue("user");
    
    $.ajax({
        url: '/api/users/' + userId + '/counts'
    }).done(function (data, textStatus, jqXHR) {
        console.log('loaded data');
        
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
            
            console.log('it\'s getting hot in here');
        }).fail(function (jqXHR, textStatus, errorThrown) {
            alert('failed to load counts');
        });
        
        $.ajax({
            url: '/api/users/' + userId + '/sync'
        }).done(function (data, textStatus, jqXHR) {
            console.log(data);
            var startDate = addHyphens(data.StartDate);
            var endDate = addHyphens(data.SyncedThroughDate);
            
            $("#startDate").val(startDate);
            $("#startDate").attr("min", startDate);
            $("#startDate").attr("max", endDate);
            $("#endDate").val(endDate);
            $("#endDate").attr("min", startDate);
            $("#endDate").attr("max", endDate);
        }).fail(function(){
            alert('failed to load sync status');
        });    
        
        console.log('whoa');
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
        var userId = getQuerystringParameterValue("user");
        var startDate = removeHyphen($("#startDate").val());
        console.log(startDate);
        var endDate = removeHyphen($("#endDate").val());
        console.log(endDate);

        $.ajax({
            url: '/api/users/' + userId + '/counts?startDate=' + startDate + '&endDate=' + endDate
        }).done(function (data, textStatus, jqXHR) {
            console.log('refreshed data');
            
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
        }).fail(function (jqXHR, textStatus, errorThrown) {
            alert('failed to load counts');
        });    
    }
    
    $(document).ready(function () {
        heatUp();
    });