$(function () {
  var $lifecycle = $('body>.main>.workspace>.lifecycle')

  var $form = $lifecycle.find('form[name=lifecycle]')
    .on('rx', (e, data) => {
      // console.log('message recieved', data)
      var $vendor
      var $items = $(e.target).find('>ul.lifecycle')
      $items.find('li.attr[name=id]>.input>.static').text(data.id)
      $items.find('li.attr[name=location]>.input>.static').text(data.location)
      $items.find('li.attr[name=strain_cost]>.input>.static').text(data.strain_cost)
      $items.find('li.attr[name=grain_cost]>.input>.static').text(data.grain_cost)
      $items.find('li.attr[name=bulk_cost]>.input>.static').text(data.bulk_cost)
      $items.find('li.attr[name=yield]>.input>.static').text(data.yield)
      $items.find('li.attr[name=count]>.input>.static').text(data.count)
      $items.find('li.attr[name=gross]>.input>.static').text(data.gross)
      $items.find('li.attr[name=modified_date]>.input>.static').text(data.modified_date)
      $items.find('li.attr[name=create_date]>.input>.static').text(data.create_date)
      var $strain = $items.find('ul[name=strain]')
      $strain.find('li.attr[name=id]>.input>.static').text(data.strain.id)
      $strain.find('li.attr[name=species]>.input>.static').text(data.strain.species)
      $strain.find('li.attr[name=name]>.input>.static').text(data.strain.name)
      $strain.find('li.attr[name=ctime]>.input>.static').text(new Date(data.strain.create_date).toString())
      $vendor = $strain.find('ul[name=vendor]')
      $vendor.find('li.attr[name=name]>.input>.static').text(data.strain.vendor.name)
      var $grain
      $vendor = $vendor
      var $bulk
      $vendor = $vendor
    })
    .on('tx', (e) => {
      var $items = $(e.target).find('>ul.lifecycle')
      $.ajax({
        url: "/lifecycle",
        method: 'POST',
        data: {
          "id": $items.find('li.attr[name=id]>.input>input.live').val(),
          "location": $items.find('li.attr[name=location]>.input>input.live').val(),
          "strain_cost": $items.find('li.attr[name=strain_cost]>.input>input.live').val(),
          "grain_cost": $items.find('li.attr[name=grain_cost]>.input>input.live').val(),
          "bulk_cost": $items.find('li.attr[name=bulk_cost]>.input>input.live').val(),
          "yield": $items.find('li.attr[name=yield]>.input>input.live').val(),
          "count": $items.find('li.attr[name=count]>.input>input.live').val(),
          "gross": $items.find('li.attr[name=gross]>.input>input.live').val(),
          "strain": {
            "id": $items.find('ul[name=strain]>li>.input>input>.live').val()
          },
          "grain_substrate": {
            "id": $items.find('ul[name=grain]>li>.input>input>.live').val()
          },
          "bulk_substrate": {
            "id": $items.find('ul[name=bulk]>li>.input>input>.live').val()
          },
        },
        async: true,
        success: (result, status, xhr) => {
        },
        error: (xhr, status, err) => {
        },
      })
    })

  var $ndx = $lifecycle.find('.ndx>.rows')
    .on('refresh', e => {
      $rows = $(e.currentTarget)
      $.ajax({
        url: '/lifecycles',
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          var selected = $rows.find('.selected>.uuid').text()
          $rows.empty()
          result.forEach(r => {
            var $row = $('<div class="row hover" />')
              .append($('<div class=uuid />').text(r.id))
              .append($('<div class=created_at />').text(r.create_date))
              .append($('<div class=location />').text(r.location))
            if (r.id === selected) {
              $row.addClass('selected')
            }
            $rows.append($row)
          })
          if ($rows.find('.selected').length == 0) {
            $rows.find('.row').first().click()
          }
        },
        error: (xhr, status, err) => {
          console.log(xhr, status, err)
        },
      })
    })
    .on('click', '>.row', e => {
      var $row = $(e.currentTarget)
      $.ajax({
        url: '/lifecycle/' + $row.find('div.uuid').text(),
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          // console.log('triggering with result', result)
          $form.trigger('rx', result)
        },
        error: (xhr, status, err) => {
          console.log(xhr, status, err)
        },
      })
      $row
        .parent()
        .find('.row.selected')
        .removeClass('selected')
      $row.addClass('selected')
    })

  $form.find('ul>li.attr.expandable').on('click', '>.label', e => {
    console.log('clicked', $(e.currentTarget)
      .parent()
      .toggleClass('expanded'))
  })

  $lifecycle
    .on('init', e => { })
    .on('activate', e => {
      $lifecycle
        .addClass('active')
        .find('>.ndx>.rows')
        .trigger('refresh')
    })
})