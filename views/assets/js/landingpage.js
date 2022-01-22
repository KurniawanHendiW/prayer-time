// HEADER ANIMATION
window.onscroll = function() {scrollFunction()};
var element = document.getElementById("body");
function scrollFunction() {
  if (document.body.scrollTop > 400 || document.documentElement.scrollTop > 400) {
      $(".navbar").addClass("fixed-top");
      element.classList.add("header-small");
      $("body").addClass("body-top-padding");

  } else {
      $(".navbar").removeClass("fixed-top");
      element.classList.remove("header-small");
      $("body").removeClass("body-top-padding");
  }
}

var dataResult = null
// OWL-CAROUSAL
$('.owl-carousel').owlCarousel({
    items: 3,
    loop:true,
    nav:false,
    dot:true,
    autoplay: true,
    slideTransition: 'linear',
    autoplayHoverPause: true,
    responsive:{
      0:{
          items:1
      },
      600:{
          items:2
      },
      1000:{
          items:3
      }
  }
})

function copyData(){
  /* Get the text field */
  var copyText = document.getElementById("urlData");

  /* Select the text field */
  copyText.select();
  copyText.setSelectionRange(0, 99999); /* For mobile devices */

  /* Copy the text inside the text field */
  navigator.clipboard.writeText(copyText.value);
  $("#message").html("Success copy url to clipboard");
  setTimeout(() => {
    $("#message").html(window.dataResult.message)
  }, 3000);
  /* Alert the copied text */
  // alert("Copied the text: " + copyText.value);
}
function scrollDown() {
  window.scrollBy(0, 650);
}

// SCROLLSPY
$(document).ready(function() {
  $.ajaxSetup({
    contentType: "application/json; charset=utf-8"
  });
  
  $(".nav-link").click(function() {
      var t = $(this).attr("href");
      $("html, body").animate({
          scrollTop: $(t).offset().top - 75
      }, {
          duration: 1000,
      });
      $('body').scrollspy({ target: '.navbar',offset: $(t).offset().top });
      return false;
  });

  $('#city').select2({
    selectOnClose: true,
    ajax: {
      url: 'https://prayer-time-calendar.herokuapp.com/prayer-time/get-city',
      data: function (params) {
        var query = {
          name: params.term
        }
  
        // Query parameters will be ?search=[term]&type=public
        return query;
      },
      processResults: function (data) {
        // Transforms the top-level key of the response object from 'items' to 'results'
        var city = []
        data.forEach(index => {
          var data = {
            "id": index.cityCode,
            "text": index.cityName
          }
          city.push(data)
        });
        return {
          results: city
        };
      }
    }
  });

  $('#day').select2({
    maximumSelectionLength: 7
  });

  $('#sholat').select2({
    maximumSelectionLength: 8
  });

  $("#get_data").click(function(){
    
    var form = document.querySelector('#prayerForm');
    var formData = new FormData(form);
    formData.set('day', $("#day").val());
    formData.set('sholat', $("#sholat").val());
    var params = {}
    formData.forEach((value, key) => params[key] = value);
    params.day = params.day.split(',')
    params.sholat = params.sholat.split(',')
    params = JSON.stringify(params);
    console.log(params, "param")
    
    $.post("https://prayer-time-calendar.herokuapp.com/prayer-time/get-key",
    params,
    function(data, status){
      // alert("Data: " + data + "\nStatus: " + status);
      window.dataResult = data
      $(".result-url").val(data.url)
      $("#message").html(data.message)
      $(".result-group").css("display", "block");
      console.log(window.dataResult , status);
    });
  });

  
});

// AOS
AOS.init({
    offset: 120, 
    delay: 0,
    duration: 1200, 
    easing: 'ease', 
    once: true, 
    mirror: false, 
    anchorPlacement: 'top-bottom', 
    disable: "mobile"
  });

//SIDEBAR-OPEN
  $('#navbarSupportedContent').on('hidden.bs.collapse', function () {
    $("body").removeClass("sidebar-open");
  })

  $('#navbarSupportedContent').on('shown.bs.collapse', function () {
    $("body").addClass("sidebar-open");
  })


  window.onresize = function() {
    var w = window.innerWidth;
    if(w>=992) {
      $('body').removeClass('sidebar-open');
      $('#navbarSupportedContent').removeClass('show');
    }
  }