def time_data_for_week(week_data,start,employee_id)

  data = Hash.new { |h,k| h[k] = Hash.new }

  %w[M T W TH F S].zip((0..6).to_a).each do |day,offset|

    date = (start + offset.days).beginning_of_day

    data[day][:lunch_hours] = LunchTime.find(:all, conditions:
      ["employee_id = ? and day between ? and ?",
          employee_id, date, date + 1.day - 1.second] ).inject(0) { |s,r|
            s + r.duration
          }

   times = [[:sick_hours    , "Sick"    ],
            [:personal_hours, "Personal"],
            [:vacation_hours, "Vacation"],
            [:other_hours,    "Other"   ]]

   times.each do |a,b|
     data[day][a] = OtherTime.find(:all, conditions:
       ["employee_id = ? and category = '#{b}' and date between ? and ?",
         employee_id, date, date + 1.day - 1.second] ).inject(0) { |s,r|
           s + r.hours
         }
   end

    d = week_data.find { |d,_| d == date }
    next unless d

    d = d[-1]
    data[day].merge!(
      regular_hours: d.inject(0) { |s,e|
        s + (e.end_time ? (e.end_time - e.start_time) / 3600 : 0)
      } - data[day][:lunch_hours],
      start_time: d.map { |e| e.start_time }.sort[0],
        end_time: d.map { |e| e.end_time }.compact.sort[-1]
    )
  end

  sums = Hash.new(0)

  data.each do |k,v|
    [:regular_hours, :lunch_hours, :sick_hours,
     :personal_hours, :vacation_hours, :other_hours].each { |h|
       sums[h] += v[h].to_f }
   end

  Table(:day,:start_time,:end_time,:regular_hours,:lunch_hours,
        :sick_hours,:personal_hours,:vacation_hours, :other_hours) do |t|
     %w[M T W TH F S].each { |d| t << {day: d}.merge(data[d]) }
     t << []
     t << { day: "<b>Totals</b>" }.merge(sums)
  end
end
