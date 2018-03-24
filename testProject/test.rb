require 'foo'

class Bar
	def initialize
		@foo = Foo.new
	end

	def do_stuff
		do_that(@foo)
	end

	def do_foo
		@foo.results.each do |f|
			puts f.bar
		end
	end
end
