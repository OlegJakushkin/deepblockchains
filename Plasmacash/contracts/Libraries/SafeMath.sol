/** This file is part of the Plasmacash library.
*
* The plasmacash library is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The Plasmacash library is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
*
* You should have received a copy of the GNU Lesser General Public License
* along with the plasmacash library. If not, see <http://www.gnu.org/licenses/>.
*/

pragma solidity ^0.4.18;

library SafeMath {
	    function mul(uint256 a, uint256 b)
	        internal
	        pure
	        returns (uint256)
	    {
	        if (a == 0) {
	        return 0;
	        }
	        uint256 c = a * b;
	        assert(c / a == b);
	        return c;
	    }

	    function div(uint256 a, uint256 b)
	        internal
	        pure
	        returns (uint256)
	    {
	        // assert(b > 0); // Solidity automatically throws when dividing by 0
	        uint256 c = a / b;
	        // assert(a == b * c + a % b); // There is no case in which this doesn't hold
	        return c;
	    }

	    function sub(uint256 a, uint256 b)
	        internal
	        pure
	        returns (uint256)
	    {
	        assert(b <= a);
	        return a - b;
	    }

	    function add(uint256 a, uint256 b)
	        internal
	        pure
	        returns (uint256)
	    {
	        uint256 c = a + b;
	        assert(c >= a);
	        return c;
	    }
}
